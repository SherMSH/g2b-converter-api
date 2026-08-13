package service

import (
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
)

func SetCardStatusG2b(pan, expDate, newStatus, reason string) (cardInfo *d8corp.CardInfoData, err error) {
	var resp *d8corp.CommonResp

	cardInfo, err = GetCardInfo(pan, expDate)
	if err != nil {
		logger.Errorf("Error Getting CardInfo data: %v", err)
		return nil, err
	}

	req := d8corp.SetCardStatusReq{
		CardKey: d8corp.CardKey{
			Pan:        pan,
			ExpiryDate: expDate,
		},
		NewStatCode: newStatus,
		Reason:      reason,
		Force:       true,
	}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b setCardStatus REQ marshaling err: %v", err)
		return nil, fmt.Errorf("[SERVICE] D8 G2b setCardStatus REQ marshaling err")
	}
	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/setCardStatus", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b setCardStatus request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b setCardStatus resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b setCardStatus RESP marshaling err: %v", err)
		return nil, err
	}
	if resp.Status.Code != "0" {
		logger.Errorf("[SERVICE] D8 G2b setCardStatus RESP status %s", resp.Status.Code)
		return nil, fmt.Errorf("%s - %s", resp.Status.RspCode, resp.Status.Message)
	}

	if cardInfo.CardBasicInfo.StatCode == "03" && newStatus == "00" {
		GeneratePIN(pan, expDate)
	}

	return
}
