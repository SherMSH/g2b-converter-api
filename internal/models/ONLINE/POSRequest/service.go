package posrequestrq

import (
	d8corp "converterapi/internal/models/D8CORP"
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"fmt"
	"strconv"
	"strings"
)

// isDebit - операция расходует средства.
//
// Статусы «только приход» - у карты 4, у счёта 2 и 4 - отклоняют именно расход.
// По эталонным ответам партнёра на карте со статусом 4 проходят запрос баланса
// (117) и зачисление (140), а проверка карты (116), перевод (135) и платёж
// (175) отклоняются кодом 58. Поэтому проверка карты считается расходной:
// она выполняется перед списанием и по ограниченной карте не разрешена.
func isDebit(txnType utils.TxnType) bool {
	switch txnType {
	case utils.Deposit, utils.A2C, utils.H2C, utils.Balance:
		return false
	default:
		return true
	}
}

func PosReq(body *Body) (soapResp *Envelope, err error) {
	//Basic checkups
	// Тип операции проверяем до InitiateTransaction: незачем занимать ссылку
	// в процессинге под запрос, который всё равно будет отклонён
	txnType := body.SoapRq.Req.GetTxnType()
	if txnType == "" {
		logger.Errorf("PosReq error: unsupported TranCode %v", body.SoapRq.Req.TranCode)
		return nil, fmt.Errorf("PosReq error: unsupported 'TranCode' value")
	}
	// Проверка карты и баланса идёт нулевой суммой, для остальных операций она обязательна
	if txnType != utils.Accver && txnType != utils.Balance && body.SoapRq.Req.Amount <= 0. {
		logger.Errorf("PosReq error: Wrong 'Amount' field value")
		return nil, fmt.Errorf("PosReq error: Wrong 'Amount' field value")
	}
	// Статус карты и счёта проверяем до обращения к процессингу.
	//
	// По документу партнёра ответ на POSRequest определяется статусом: истёкшая
	// карта - 51, требующая обращения к эмитенту - 71 и так далее. Процессинг
	// такие операции пропускает и списывает средства, поэтому проверять после
	// авторизации поздно: деньги уже ушли, а партнёру надо вернуть отказ.
	preCard, err := service.GetCardInfo(body.SoapRq.Req.PAN, body.SoapRq.Req.GetExpDate())
	if err != nil {
		logger.Warnf("[SERVICE] POSRequest: статус карты не проверен: %v", err)
	}
	if code := statusDeclineCode(preCard, isDebit(txnType)); code != "" {
		logger.Warnf("[SERVICE] POSRequest declined by status: TWO %s, PAN %s", code, preCard.CardBasicInfo.Lkey.MaskedPan)
		return statusDeclineResponse(body, code, preCard, isDebit(txnType)), nil
	}

	ectxNum, err := service.InitiateTransaction()
	if err != nil {
		logger.Errorf("POS req {InitiateTransaction} error: %v", err)
		return nil, err
	}
	logger.Infof("PosReq info: %+v", body.SoapRq.Req)
	trn, err := service.AuthorizeTransaction(body.SoapRq.Req, *ectxNum)
	if err != nil {
		logger.Errorf("POS req {AuthorizeTransaction} error: %v", err)
		return nil, err
	}
	if trn != nil {
		body.SoapRq.Req.ThisTranId = fmt.Sprintf("%v", trn.TransactionResponse.TlId)
		body.SoapRq.Req.RespCode = trn.TransactionResponse.RspCode
		body.SoapRq.ApprovalCode = trn.TransactionResponse.ApprovalCode
	}
	// Операция не принята к обработке: это состояние транзакции, а не код ответа.
	// Раньше здесь сравнивался rspCode, у которого 17 означает Bad PIN, - отказ по
	// неверному PIN превращался в ошибку сервиса вместо ответа с кодом причины.
	if strconv.Itoa(trn.TransactionResponse.TxStatus) == string(utils.AdviceLogNotProceed) {
		logger.Errorf("bad response tx status {Skipped}")
		return nil, fmt.Errorf("bad response tx status {Skipped}")
	}

	// Собираем SoapEnvelope Response
	var (
		cardInfo       *d8corp.CardInfoData
		accnum, cvok   string
		avlbal, blkamt float64
	)
	trnDetails, err := service.GetTransactionDetailsG2b(fmt.Sprintf("%d", trn.TransactionResponse.TlId), trn.TransactionResponse.EcTxRefno)
	if err != nil {
		logger.Errorf("[SERVICE] POSRequest error getting trn details")
	}
	// Данные карты запрашиваем заново: партнёру нужны балансы на момент после
	// авторизации, а preCard получена до неё.
	if trnDetails != nil {
		cardInfo, err = service.GetCardInfo(trn.Lkey.Pan, trnDetails.Details.DateExp)
		if err != nil {
			logger.Errorf("[SERVICE] POSRequest error getting card info: %v", err)
		}
	}
	if cardInfo != nil && len(cardInfo.CardAccounts) != 0 {
		accnum = cardInfo.CardAccounts[0].AccountNumber
		avlbal = cardInfo.CardAccounts[0].AvlBal
		blkamt = cardInfo.CardAccounts[0].BlkAmt
	}

	// Код ответа авторизатора партнёр ждёт в кодировке TWO, а процессинг отвечает
	// парой code/rspcode - переводим по справочнику.
	authRespCode := utils.AuthRespCode(trn.TransactionResponse.ActionCode, trn.TransactionResponse.RspCode)

	if !utils.IsApproved(trn.TransactionResponse.ActionCode) {
		authRespCode = refineDeclineCode(authRespCode, cardInfo, isDebit(txnType))
	}

	// Причина отказа обязана быть в каждом неуспешном ответе. Процессинг её
	// заполняет не всегда, поэтому при пустом сообщении собираем текст сами -
	// по статусу карты или счёта, а если и он ни при чём, по коду ответа.
	declineReason := trn.DeclineReason
	if !utils.IsApproved(trn.TransactionResponse.ActionCode) && declineReason == "" {
		declineReason = buildDeclineReason(authRespCode, cardInfo, isDebit(txnType))
	}

	if !utils.IsApproved(trn.TransactionResponse.ActionCode) {
		logger.Warnf("[SERVICE] POSRequest declined: D8 %s/%s -> TWO %s, tlId %d",
			trn.TransactionResponse.ActionCode, trn.TransactionResponse.RspCode,
			authRespCode, trn.TransactionResponse.TlId)
	}

	// Детали и данные карты могли не прийти: собираем ответ из того, что есть
	// Валюты в POS-ответе партнёр ждёт числовыми кодами (972), в отличие от
	// выписок, где используется буквенный код
	var accountCurrency, balanceCurrency, billCurrency, toAcct string
	if cardInfo != nil {
		accountCurrency = cardInfo.CardBasicInfo.Currcode
	}
	if trnDetails != nil {
		balanceCurrency = trnDetails.Details.TxnCurrency
		billCurrency = trnDetails.Details.Curbill
		toAcct = trnDetails.Details.DestinationAccountType
	}

	cvok = "-1"
	if txnType == utils.Accver || txnType == utils.Balance {
		cvok = "1"
	}

	// Балансы отдаём только по одобренной операции. Отказ означает, что карта
	// или счёт не допущены к работе, и раскрывать по ним остатки не следует -
	// особенно при запросе баланса, ради которого операция и делалась.
	availBalance, ledgerBalance := "", ""
	if utils.IsApproved(trn.TransactionResponse.ActionCode) {
		availBalance = fmt.Sprintf("%.2f", avlbal)
		ledgerBalance = fmt.Sprintf("%.2f", avlbal+blkamt)
	}

	soapResp = &Envelope{
		XmlnsS:  "http://www.w3.org/2003/05/soap-envelope",
		XmlnsM1: "http://schemas.compassplus.com/two/1.0/fimi.xsd",
		XmlnsM0: "http://schemas.compassplus.com/two/1.0/fimi_types.xsd",
		Body: RespBody{
			POSRequestRp: POSRequestRp{
				Response: Response{
					Product:      body.SoapRq.Req.Product,
					ResponseAttr: "1",
					TranId:       utils.GenerateTimestampID(),
					Ver:          "1.0",
					Echo:         body.SoapRq.Req.Echo,

					AccountCurrency:      accountCurrency,
					ApprovalCode:         trn.TransactionResponse.ApprovalCode,
					AuthRespCode:         authRespCode,
					AuthRespCodeCategory: "0",
					AvailBalance:         availBalance,
					BalanceCurrency:      balanceCurrency,
					BonusDebt:            "0",
					CVxOK:                cvok,
					Currency:             billCurrency,
					DeclineReason:        declineReason,
					Fee:                  "0",
					FromAcct:             accnum,
					IssuerFee:            "0",
					LedgerBalance:        ledgerBalance,
					MaskBalances:         "0",
					RelatedTran:          RelatedTran{},
					ThisTranId:           fmt.Sprintf("%d", trn.TransactionResponse.TlId),
					ToAcct:               toAcct,
				},
			},
		},
	}
	return soapResp, nil
}

// buildDeclineReason собирает текст причины отказа, когда процессинг прислал
// пустое сообщение.
//
// Формат повторяет тот, что партнёр видит от TWO:
// "Response for card status 'Lost' in authorization scheme #1. Card #505827******0016".
// Если статус карты и счёта операцию не запрещают, причина неизвестна - тогда
// отдаём хотя бы код ответа, но не оставляем поле пустым.
func buildDeclineReason(authRespCode string, cardInfo *d8corp.CardInfoData, isDebit bool) string {
	if cardInfo != nil {
		cardStatus := utils.CardStatuses[cardInfo.CardBasicInfo.StatCode]
		if utils.CardStatusRespCode(cardStatus, isDebit) != "" {
			return fmt.Sprintf("Response for card status '%s' in authorization scheme #1. Card #%s",
				utils.CardStatusNames[cardStatus], maskPan(cardInfo.CardBasicInfo.Lkey.Pan))
		}

		if len(cardInfo.CardAccounts) != 0 {
			acct := cardInfo.CardAccounts[0]
			acctStatus := utils.AccountStatuses[acct.StatCode]
			switch {
			// Неактивный и закрытый счёт партнёр видит отдельной формулировкой
			case acctStatus == "0" || acctStatus == "9":
				return fmt.Sprintf("Account #%s is inactive or closed (GetPrimaryAccount)", acct.AccountNumber)
			case utils.AccountStatusRespCode(acctStatus, isDebit) != "":
				return fmt.Sprintf("Status '%s' of account #%s is not appropriate for this transaction (GetPrimaryAccount)",
					acctStatus, acct.AccountNumber)
			}
		}
	}
	return fmt.Sprintf("Transaction declined with response code %s", authRespCode)
}

// maskPan приводит номер карты к виду 976249******4049 - так он выглядит в
// текстах причин, которые партнёр получает от TWO.
func maskPan(pan string) string {
	if len(pan) < 10 {
		return pan
	}
	return pan[:6] + strings.Repeat("*", len(pan)-10) + pan[len(pan)-4:]
}

// statusDeclineCode возвращает код отказа, если статус карты или счёта
// запрещает операцию. Пустая строка означает, что статусы операции не мешают.
func statusDeclineCode(cardInfo *d8corp.CardInfoData, isDebit bool) string {
	if cardInfo == nil {
		return ""
	}
	if code := utils.CardStatusRespCode(utils.CardStatuses[cardInfo.CardBasicInfo.StatCode], isDebit); code != "" {
		return code
	}
	if len(cardInfo.CardAccounts) != 0 {
		if code := utils.AccountStatusRespCode(utils.AccountStatuses[cardInfo.CardAccounts[0].StatCode], isDebit); code != "" {
			return code
		}
	}
	return ""
}

// statusDeclineResponse собирает отказ по статусу карты или счёта.
//
// Транзакция в процессинг не отправлялась, поэтому в ответе нет ни номера
// операции, ни кода авторизации, ни балансов - только причина отказа.
func statusDeclineResponse(body *Body, authRespCode string, cardInfo *d8corp.CardInfoData, isDebit bool) *Envelope {
	// Транзакции не было, поэтому валюту берём из запроса партнёра, а при её
	// отсутствии - из данных карты
	currency := body.SoapRq.Req.Currency
	if currency == "" && cardInfo != nil {
		currency = cardInfo.CardBasicInfo.Currcode
	}

	return &Envelope{
		XmlnsS:  "http://www.w3.org/2003/05/soap-envelope",
		XmlnsM1: "http://schemas.compassplus.com/two/1.0/fimi.xsd",
		XmlnsM0: "http://schemas.compassplus.com/two/1.0/fimi_types.xsd",
		Body: RespBody{
			POSRequestRp: POSRequestRp{
				Response: Response{
					Product:      body.SoapRq.Req.Product,
					ResponseAttr: "1",
					TranId:       utils.GenerateTimestampID(),
					Ver:          "1.0",
					Echo:         body.SoapRq.Req.Echo,

					AuthRespCode:         authRespCode,
					AuthRespCodeCategory: "0",
					BalanceCurrency:      currency,
					BonusDebt:            "0",
					CVxOK:                "-1",
					Currency:             currency,
					DeclineReason:        buildDeclineReason(authRespCode, cardInfo, isDebit),
					Fee:                  "0",
					MaskBalances:         "0",
					RelatedTran:          RelatedTran{},
					ToAcct:               body.SoapRq.Req.ToAcctType,
				},
			},
		},
	}
}

// refineDeclineCode уточняет код отказа по статусу карты и счёта.
//
// Процессинг часто отвечает общим отказом - "Do not honour" или неизвестным
// кодом, - и партнёр получает 50 или 68 вместо настоящей причины. По таблице
// партнёра код определяется статусом: скомпрометирована - 75, потеряна - 40,
// украдена - 41 и так далее. Поэтому статус карты важнее общего кода.
//
// Конкретные коды процессинга (недостаточно средств, неверный PIN) не трогаем:
// они точнее любого статуса.
func refineDeclineCode(authRespCode string, cardInfo *d8corp.CardInfoData, isDebit bool) string {
	if cardInfo == nil || !isGenericDecline(authRespCode) {
		return authRespCode
	}
	if code := statusDeclineCode(cardInfo, isDebit); code != "" {
		return code
	}
	return authRespCode
}

// isGenericDecline - отказ без конкретной причины: "несанкционированное
// использование" и "отказ внешнего хоста". Оба ничего не говорят о том, что
// именно не так с картой.
func isGenericDecline(authRespCode string) bool {
	return authRespCode == "50" || authRespCode == utils.TwoExternalDecline
}
