package service

import (
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
)

func GetCardInfo(pan, expdate string) (resp *d8corp.CommonResp, err error) {
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
		logger.Errorf("[SERVICE] D8 G2b authorizeTransaction request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b authorizeTransaction resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b authorizeTransaction RESP marshaling err: %v", err)
		return nil, err
	}
	return resp, nil
}
