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
	// Данные карты нужны и для балансов, и для проверки статуса после авторизации
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

	declinedByStatus := ""
	if utils.IsApproved(trn.TransactionResponse.ActionCode) {
		// Процессинг пропускает операции по картам со статусами «истёк срок» и
		// «нужен запрос к эмитенту», хотя по документу партнёра они должны
		// отклоняться. Отменяем такую операцию и отдаём отказ: номер транзакции
		// при этом настоящий, как в эталонных ответах.
		declinedByStatus = statusDeclineCode(cardInfo, isDebit(txnType))
		if declinedByStatus != "" {
			if err := reverseApproved(trn, body.SoapRq.Req.Amount, body.SoapRq.Req.Currency); err != nil {
				// Средства списаны, отменить не удалось. Сообщить об отказе
				// значило бы разойтись с процессингом: у клиента деньги ушли, а
				// партнёр считал бы операцию непрошедшей.
				logger.Errorf("[SERVICE] POSRequest: операция %d одобрена процессингом вопреки статусу карты, отменить не удалось: %v",
					trn.TransactionResponse.TlId, err)
				declinedByStatus = ""
			}
		}
	}

	if declinedByStatus != "" {
		authRespCode = declinedByStatus
	} else if !utils.IsApproved(trn.TransactionResponse.ActionCode) {
		authRespCode = refineDeclineCode(authRespCode, cardInfo, isDebit(txnType))
	}

	// Причина отказа обязана быть в каждом неуспешном ответе.
	//
	// Когда операцию не пропускает статус карты или счёта, текст собираем сами в
	// формате TWO: процессинг пишет по-своему и для одного и того же статуса
	// по-разному - на платёж "PIN tries exceeded has been set on card status", на
	// перевод "AFT Rejected". Партнёр же ждёт "Response for card status ...".
	// Сообщение процессинга оставляем только там, где оно действительно несёт
	// причину: на нехватку средств оно приходит то коротким "Insufficient
	// funds !", то дежурным "AFT Rejected" на переводе, а партнёр в обоих
	// случаях ждёт разбор по слагаемым остатка.
	declined := declinedByStatus != "" || !utils.IsApproved(trn.TransactionResponse.ActionCode)
	declineReason := trn.DeclineReason
	if declined {
		byStatus := statusDeclineCode(cardInfo, isDebit(txnType)) != ""
		if byStatus || authRespCode == utils.TwoInsufficientFunds || isPlaceholderReason(declineReason) {
			declineReason = buildDeclineReason(authRespCode, cardInfo, isDebit(txnType), body.SoapRq.Req.Amount)
		}
	}

	if declined {
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
	if !declined {
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
					RelatedTran:          relatedTran(trnDetails),
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
func buildDeclineReason(authRespCode string, cardInfo *d8corp.CardInfoData, isDebit bool, amount float64) string {
	if cardInfo != nil {
		cardStatus := utils.CardStatuses[cardInfo.CardBasicInfo.StatCode]
		if utils.CardStatusRespCode(cardStatus, isDebit) != "" {
			return fmt.Sprintf("Response for card status '%s' in authorization scheme #1. Card #%s",
				utils.CardStatusNames[cardStatus], maskPan(cardInfo.CardBasicInfo.Lkey.Pan))
		}

		if len(cardInfo.CardAccounts) != 0 {
			acct := cardInfo.CardAccounts[0]
			acctStatus := utils.AccountStatuses[acct.StatCode]

			// Нехватку средств партнёр ждёт с разбором по слагаемым остатка.
			// Поля, которых процессинг не отдаёт, показываем нулями - в его
			// эталонных ответах они тоже нулевые.
			if authRespCode == utils.TwoInsufficientFunds {
				return fmt.Sprintf("Insufficient funds on account #%s. Tranx amount=%.2f is more than "+
					"AcctEffectiveBalance=%.2f (AvailBalance=%.2f, OverdraftLimit=%.2f, Bonus=0.00, "+
					"TmpOverdraft=0.00, EMVOfflineHold=0.00, Protected Amount=0.00)",
					acct.AccountNumber, amount, acct.AvlBal+acct.Crlimit, acct.AvlBal, acct.Crlimit)
			}

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

// placeholderReasons - дежурные сообщения процессинга, которые причину не
// несут. Плечи перевода (AFT - списание, OCT - зачисление) отклоняются именно
// с ними, из-за чего партнёр видел "AFT Rejected" там, где ждал разбор отказа.
var placeholderReasons = map[string]bool{
	"":              true,
	"aft rejected":  true,
	"oct rejected":  true,
	"txn rejected":  true,
	"trxn rejected": true,
}

// isPlaceholderReason сообщает, что текст процессинга бесполезен и причину
// нужно собрать самим.
func isPlaceholderReason(reason string) bool {
	return placeholderReasons[strings.ToLower(strings.TrimSpace(reason))]
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

// relatedTran заполняет блок связанных операций.
//
// Процессинг возвращает связи в transactionGroups: тип 2 связывает операцию с
// её реверсом, тип 5 - части P2P-перевода, тип 4 - преавторизацию с
// завершением. Партнёр ждёт их в RelatedTran с тем же смыслом, поэтому
// переносим как есть - идентификаторы настоящие, из ответа процессинга.
func relatedTran(details *d8corp.Transaction) RelatedTran {
	if details == nil {
		return RelatedTran{}
	}

	rows := make([]Rows, 0)
	for _, group := range details.Details.TransactionGroups {
		for _, related := range group.Transactions {
			if related.TlId == details.Details.TlId {
				continue // сама операция, а не связанная с ней
			}
			rows = append(rows, Rows{
				RelatedTranId:       strconv.Itoa(related.TlId),
				RelatedTranCode:     utils.TranCode(related.TxnCode),
				RelatedAuthRespCode: utils.AuthRespCode(related.ActionCode, related.RspCode),
			})
		}
	}
	if len(rows) == 0 {
		return RelatedTran{}
	}
	return RelatedTran{Rows: rows}
}

// reverseApproved отменяет операцию, которую процессинг одобрил вопреки статусу
// карты или счёта.
//
// Реверсу нужна собственная ссылка, поэтому сначала InitiateTransaction (8.6).
// Операции с нулевой суммой - проверка карты и запрос баланса - средств не
// двигают, отменять там нечего.
func reverseApproved(trn *d8corp.TrnData, amount float64, currency string) error {
	if amount <= 0 {
		return nil
	}
	ecTxRefNo, err := service.InitiateTransaction()
	if err != nil {
		return fmt.Errorf("реверс не инициирован: %w", err)
	}
	if currency == "" {
		currency = utils.TJSCurrency
	}
	_, err = service.ReverseTransaction(*ecTxRefNo, trn.TransactionResponse.EcTxRefno, amount, currency, reversalReasonStatus)
	return err
}

// reversalReasonStatus - причина отмены: операция не соответствует условиям
// проведения. Коды причин - Appendix C спецификации D8.
const reversalReasonStatus = 4000

// refineDeclineCode определяет код отказа по статусу карты и счёта.
//
// В таблицах партнёра код однозначно задан статусом, и подпись под ними прямо
// говорит, что это ответы на запросы POSRequest: ограничена - 58,
// скомпрометирована - 75, потеряна - 40. Формулировка причины отказа
// ("in authorization scheme #1") тоже показывает, что статус проверяется первым.
// Поэтому статус важнее кода процессинга: на карте со статусом 4 процессинг
// отвечает 1/06 "PIN tries exceeded", что дало бы 62 вместо ожидаемого 58.
//
// Код процессинга остаётся, когда статусы операцию не запрещают - тогда причина
// в самой операции: нехватка средств, неверный PIN, превышение лимита.
func refineDeclineCode(authRespCode string, cardInfo *d8corp.CardInfoData, isDebit bool) string {
	if code := statusDeclineCode(cardInfo, isDebit); code != "" {
		return code
	}
	return authRespCode
}
