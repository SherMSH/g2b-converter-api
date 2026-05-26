package posrequestrq

import (
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"fmt"
)

func PosReq(body *Body) error {
	//Basic checkups
	if body.SoapRq.Req.Amount <= 0. {
		logger.Errorf("PosReq error: Wrong 'Amount' field value")
		return fmt.Errorf("PosReq error: Wrong 'Amount' field value")
	}
	if len(body.SoapRq.Req.ToAccount) == 0 {
		logger.Errorf("PosReq error: Mandatory field 'ToAccount' is missing")
		return fmt.Errorf("Mandatory field 'ToAccount' is missing")
	}

	ectxNum, err := service.InitiateTransaction()
	if err != nil {
		logger.Errorf("POS req {InitiateTransaction} error: %v", err)
		return err
	}
	logger.Infof("PosReq info: %+v", body.SoapRq.Req)
	trn, atr, err := service.AuthorizeTransaction(body.SoapRq.Req, *ectxNum)
	if err != nil {
		logger.Errorf("POS req {AuthorizeTransaction} error: %v", err)
		return err
	}

	if atr.Status.Code != "0" {
		logger.Errorf("bad response status code: %+v", atr.Status)
		return fmt.Errorf("%s", atr.Status.Message)
	}
	if trn != nil {
		body.SoapRq.Req.ThisTranId = fmt.Sprintf("%v", trn.TransactionResponse.TlId)
		body.SoapRq.Req.RespCode = trn.TransactionResponse.RspCode
		body.SoapRq.ApprovalCode = trn.ApprovalCode
	}
	if trn.TransactionResponse.RspCode == string(utils.AdviceLogNotProceed) {
		logger.Errorf("bad response tx status {Skipped}")
		return fmt.Errorf("bad response tx status {Skipped}")
	}
	return nil
}
