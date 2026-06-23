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
	if body.SoapRq.Req.Amount <= 0. {
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
		body.SoapRq.ApprovalCode = trn.ApprovalCode
	}
	if trn.TransactionResponse.RspCode == string(utils.AdviceLogNotProceed) {
		logger.Errorf("bad response tx status {Skipped}")
		return nil, fmt.Errorf("bad response tx status {Skipped}")
	}

	// Собираем SoapEnvelope Response
	var (
		cardInfo       *d8corp.CardInfoData
		accnum         string
		avlbal, blkamt float64
	)
	trnDetails, err := service.GetTransactionDetailsG2b(fmt.Sprintf("%d", trn.TransactionResponse.TlId), trn.TransactionResponse.EcTxRefno)
	if err != nil {
		logger.Errorf("[SERVICE] POSRequest error getting trn details")
	}
	if trnDetails != nil {
		cardInfo, _ = service.GetCardInfo(trn.Lkey.Pan, trnDetails.Details.DateExp)
		if len(cardInfo.CardAccounts) != 0 {
			accnum = cardInfo.CardAccounts[0].AccountNumber
			avlbal = cardInfo.CardAccounts[0].AvlBal
			blkamt = cardInfo.CardAccounts[0].BlkAmt
		}
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

					AccountCurrency:      utils.Currencies[cardInfo.CardBasicInfo.Currcode],
					ApprovalCode:         trn.ApprovalCode,
					AuthRespCode:         body.SoapRq.Req.RespCode,
					AuthRespCodeCategory: "0",
					AvailBalance:         fmt.Sprintf("%.2f", avlbal),
					BalanceCurrency:      utils.Currencies[trnDetails.Details.TxnCurrency],
					BonusDebt:            "0",
					CVxOK:                "-1",
					Currency:             utils.Currencies[trnDetails.Details.Curbill],
					Fee:                  "",
					FromAcct:             accnum,
					IssuerFee:            "",
					LedgerBalance:        fmt.Sprintf("%.2f", avlbal+blkamt),
					MaskBalances:         "0",
					RelatedTran:          RelatedTran{},
					ThisTranId:           body.SoapRq.Req.ThisTranId,
					ToAcct:               trnDetails.Details.DestinationAccountType,
				},
			},
		},
	}
	return soapResp, nil
}
