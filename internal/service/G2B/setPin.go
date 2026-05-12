package service

import (
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
)

func SetPinG2b(pan, expDate string) error {
	var resp *d8corp.CommonResp
	req := d8corp.SetPinReq{
		CardKey: d8corp.CardKey{
			Pan:        pan,
			ExpiryDate: expDate,
		},
		PinKeyUnderRSA: "",
		PinBlock:       "",
		PinBlockType:   "0",
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b setPIN REQ marshaling err: %v", err)
		return fmt.Errorf("[SERVICE] D8 G2b setPIN REQ marshaling err")
	}
	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/setPIN", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b setPIN request sending err: %v", err)
		return err
	}
	logger.Infof("[SERVICE] D8 G2b setPIN resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b setPIN RESP marshaling err: %v", err)
		return err
	}
	if resp.Status.Code != "0" {
		logger.Errorf("[SERVICE] D8 G2b setPIN RESP status %s", resp.Status.Code)
		return fmt.Errorf("%s", resp.Status.RspCode)
	}
	// err = json.Unmarshal(resp.Data, &cardInfo)
	// if err != nil {
	// 	logger.Errorf("[SERVICE] D8 G2b setPIN RESP marshaling err: %v", err)
	// 	return err
	// }
	// if cardInfo == nil {
	// 	logger.Errorf("[SERVICE] D8 G2b setPIN RESP is empty")
	// 	return fmt.Errorf("no data")
	// }
	return nil
}
