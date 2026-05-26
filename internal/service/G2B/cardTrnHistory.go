package service

import (
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
)

func GetCardTransactionHistory(pan, from, to string, size int) (transactions *d8corp.CardInfoData, err error) {
	var resp *d8corp.CommonResp
	req := d8corp.GetCardTrnHistoryReq{
		CardKey: d8corp.CardKey{
			Pan: pan,
		},
		DateLocalFrom: from,
		DateLocalTo:   to,
		PagingParams: d8corp.PagingParams{
			Size:            size,
			LastRetrievedId: -1,
		},
	}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b GetCardInfo REQ marshaling err: %v", err)
		return nil, fmt.Errorf("[SERVICE] D8 G2b GetCardInfo REQ marshaling err")
	}

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/kernel/1.0/getCardTransactionHistory", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b GetCardInfo request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b GetCardInfo resp status: %v, body: %v", status, string(data))
	err = json.Unmarshal(data, &resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b GetCardInfo RESP marshaling err: %v", err)
		return nil, err
	}
	if resp.Status.Code != "0" {
		logger.Errorf("[SERVICE] D8 G2b GetCardInfo RESP status %s", resp.Status.Code)
		return nil, fmt.Errorf("%s - %s", resp.Status.RspCode, resp.Status.Message)
	}

	err = json.Unmarshal(resp.Data, &transactions)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b GetCardInfo marshaling err: %v", err)
		return nil, err
	}
	if transactions == nil {
		logger.Errorf("[SERVICE] D8 G2b GetCardInfo RESP is empty")
		return nil, fmt.Errorf("no data")
	}
	return transactions, nil
}
