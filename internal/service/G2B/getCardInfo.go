package service

import (
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
)

func GetCardInfo(pan, expdate string) (cardInfo *d8corp.CardInfoData, err error) {
	var resp *d8corp.CommonResp
	req := d8corp.GetCardInfoReq{
		CardKey: d8corp.CardKey{
			Pan:        pan,
			ExpiryDate: expdate,
		},
		ReqCardBasicInfo:        true,
		ReqCardAccounts:         true,
		ReqCardLimits:           true,
		ReqCardAccountLimits:    true,
		ReqCardAuthRestrictions: true,
		ReqCardTransactions:     true,
		CardTransactionCount:    3,
	}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b GetCardInfo REQ marshaling err: %v", err)
		return nil, fmt.Errorf("[SERVICE] D8 G2b GetCardInfo REQ marshaling err")
	}
	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/getCardInfo", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b GetCardInfo request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b GetCardInfo resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b GetCardInfo RESP marshaling err: %v", err)
		return nil, err
	}
	if resp.Status.Code != "0" {
		logger.Errorf("[SERVICE] D8 G2b GetCardInfo RESP status %s", resp.Status.Code)
		return nil, fmt.Errorf("%s", resp.Status.RspCode)
	}
	err = json.Unmarshal(resp.Data, &cardInfo)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b GetCardInfo RESP marshaling err: %v", err)
		return nil, err
	}
	if cardInfo == nil {
		logger.Errorf("[SERVICE] D8 G2b GetCardInfo RESP is empty")
		return nil, fmt.Errorf("no data")
	}
	return cardInfo, nil
}
