package service

import (
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
)

func SetCardStatusG2b(pan, expDate, newStatus, reason string) (err error) {
	var resp *d8corp.CommonResp

	req := d8corp.SetCardStatusReq{
		CardKey: d8corp.CardKey{
			Pan:        pan,
			ExpiryDate: expDate,
		},
		NewStatCode: newStatus,
		Reason:      reason,
	}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b setCardStatus REQ marshaling err: %v", err)
		return fmt.Errorf("[SERVICE] D8 G2b setCardStatus REQ marshaling err")
	}
	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/setCardStatus", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b setCardStatus request sending err: %v", err)
		return err
	}
	logger.Infof("[SERVICE] D8 G2b setCardStatus resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b setCardStatus RESP marshaling err: %v", err)
		return err
	}
	if resp.Status.Code != "0" {
		logger.Errorf("[SERVICE] D8 G2b setCardStatus RESP status %s", resp.Status.Code)
		return fmt.Errorf("%s - %s", resp.Status.RspCode, resp.Status.Message)
	}
	return
}
