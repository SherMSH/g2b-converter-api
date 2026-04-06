package posrequestrq

import (
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"fmt"
)

func PosReq(body *Body) error {
	//Basic checkups
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
	body.SoapRq.Req.ThisTranId = fmt.Sprintf("%v", trn.TransactionResponse.TlId)
	body.SoapRq.Req.RespCode = trn.TransactionResponse.RspCode
	body.SoapRq.ApprovalCode = trn.ApprovalCode

	if atr.Status.Code != "0" {
		return fmt.Errorf("bad response status code: %+v", atr)
	}
	if trn.TransactionResponse.RspCode == string(utils.AdviceLogNotProceed) {
		return fmt.Errorf("bad response tx status {Skipped}")
	}
	return nil
}
