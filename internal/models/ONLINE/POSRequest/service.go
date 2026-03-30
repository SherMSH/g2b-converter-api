package posrequestrq

import (
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"fmt"
)

func PosReq(body Body) error {
	ectxNum, err := service.InitiateTransaction()
	if err != nil {
		logger.Errorf("POS req {InitiateTransaction} error: %v", err)
		return err
	}
	atr, err := service.AuthorizeTransaction(body.SoapRq.Req, *ectxNum)
	if err != nil {
		logger.Errorf("POS req {AuthorizeTransaction} error: %v", err)
		return err
	}
	if atr.Status.Code == string(utils.AdviceLogNotProceed) {
		return fmt.Errorf("bad response status {Skipped}")
	}
	return nil
}
