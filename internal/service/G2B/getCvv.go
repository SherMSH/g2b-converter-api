package service

import (
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
)

func GetCVVG2b(pan, expdate string) (cvvData *d8corp.CVVData, err error) {
	var resp *d8corp.CommonResp
	req := d8corp.GetCVVReq{
		CardKey: d8corp.CardKey{
			Pan:        pan,
			ExpiryDate: expdate,
		},
	}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b GetCVVG2b REQ marshaling err: %v", err)
		return nil, fmt.Errorf("[SERVICE] D8 G2b GetCVV2 REQ marshaling err")
	}
	logger.Infof("[SERVICE] D8 G2b GetCVV2 REQ %v", string(jsonReq))
	data, status, err := utils.SendRequest("POST", "http://d8-prod-proc-web1.humo.lab"+"/xapi/miss/1.0/getCVV2", jsonReq, utils.D8HeadersMap) // "перенапрвление в прод (HSM)"
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b GetCVV2 request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b GetCVV2 resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b GetCVV2 RESP marshaling err: %v", err)
		return nil, err
	}
	if resp.Status.Code != "0" {
		logger.Errorf("[SERVICE] D8 G2b GetCVV2 RESP status %s", resp.Status.Code)
		return nil, fmt.Errorf("%s - %s", resp.Status.RspCode, resp.Status.Message)
	}

	err = json.Unmarshal(resp.Data, &cvvData)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b GetCardInfo marshaling err: %v", err)
		return nil, err
	}
	if cvvData == nil {
		logger.Errorf("[SERVICE] D8 G2b GetCardInfo RESP is empty")
		return nil, fmt.Errorf("no data")
	}
	return cvvData, nil
}
