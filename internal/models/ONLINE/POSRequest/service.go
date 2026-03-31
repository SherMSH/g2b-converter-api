package posrequestrq

import (
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"fmt"
)

func PosReq(body *Body) error {
	ectxNum, err := service.InitiateTransaction()
	if err != nil {
		logger.Errorf("POS req {InitiateTransaction} error: %v", err)
		return err
	}
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
