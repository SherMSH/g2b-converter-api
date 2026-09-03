package posrequestrq

import (
	d8corp "converterapi/internal/models/D8CORP"
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"fmt"
)

func PosReq(body *Body) (soapResp *Envelope, err error) {
	//Basic checkups
	// Тип операции проверяем до InitiateTransaction: незачем занимать ссылку
	// в процессинге под запрос, который всё равно будет отклонён
	txnType := body.SoapRq.Req.GetTxnType()
	if txnType == "" {
		logger.Errorf("PosReq error: unsupported TranCode %v", body.SoapRq.Req.TranCode)
		return nil, fmt.Errorf("PosReq error: unsupported 'TranCode' value")
	}
	// Проверка карты идёт нулевой суммой, для остальных операций она обязательна
	if txnType != utils.Accver && body.SoapRq.Req.Amount <= 0. {
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
	if trn.TransactionResponse.RspCode == string(utils.AdviceLogNotProceed) {
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
	if trnDetails != nil {
		cardInfo, err = service.GetCardInfo(trn.Lkey.Pan, trnDetails.Details.DateExp)
		if err != nil {
			logger.Errorf("[SERVICE] POSRequest error getting card info: %v", err)
		}
		if cardInfo != nil && len(cardInfo.CardAccounts) != 0 {
			accnum = cardInfo.CardAccounts[0].AccountNumber
			avlbal = cardInfo.CardAccounts[0].AvlBal
			blkamt = cardInfo.CardAccounts[0].BlkAmt
		}
	}

	// Детали и данные карты могли не прийти: собираем ответ из того, что есть
	var accountCurrency, balanceCurrency, billCurrency, toAcct string
	if cardInfo != nil {
		accountCurrency = utils.Currency(cardInfo.CardBasicInfo.Currcode)
	}
	if trnDetails != nil {
		balanceCurrency = utils.Currency(trnDetails.Details.TxnCurrency)
		billCurrency = utils.Currency(trnDetails.Details.Curbill)
		toAcct = trnDetails.Details.DestinationAccountType
	}

	cvok = "-1"
	if txnType == utils.Accver || txnType == utils.Balance {
		cvok = "1"
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
					AuthRespCode:         body.SoapRq.Req.RespCode,
					AuthRespCodeCategory: "0",
					AvailBalance:         fmt.Sprintf("%.2f", avlbal),
					BalanceCurrency:      balanceCurrency,
					BonusDebt:            "0",
					CVxOK:                cvok,
					Currency:             billCurrency,
					Fee:                  "",
					FromAcct:             accnum,
					IssuerFee:            "",
					LedgerBalance:        fmt.Sprintf("%.2f", avlbal+blkamt),
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
