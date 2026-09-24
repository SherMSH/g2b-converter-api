package service

import (
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
)

func ResetCardPINTriesG2b(pan, expDate string) (err error) {
	var resp *d8corp.CommonResp

	req := d8corp.GetCardInfoReq{
		CardKey: d8corp.CardKey{
			Pan:        pan,
			ExpiryDate: expDate,
		},
	}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b resetCardPINTries REQ marshaling err: %v", err)
		return fmt.Errorf("[SERVICE] D8 G2b resetCardPINTries REQ marshaling err")
	}
	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/resetCardPINTries", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b resetCardPINTries request sending err: %v", err)
		return err
	}
	logger.Infof("[SERVICE] D8 G2b resetCardPINTries resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b resetCardPINTries RESP marshaling err: %v", err)
		return err
	}
	if resp.Status.Code != "0" {
		logger.Errorf("[SERVICE] D8 G2b resetCardPINTries RESP status %s", resp.Status.Code)
		return fmt.Errorf("%s - %s", resp.Status.RspCode, resp.Status.Message)
	}
	return
}

// GeneratePIN поручает процессингу сгенерировать новый PIN карты (xmiss/generatePIN, 7.7).
//
// Само значение наружу не отдаётся - ни нам, ни вызывающей стороне: по
// спецификации "Due to security reasons PIN is not returned in the response".
// До держателя карты новый PIN доходит оповещением процессинга, поэтому у карты
// должен быть заведён контракт SMS.
func GeneratePIN(pan, expDate string) error {
	var resp *d8corp.CommonResp
	req := d8corp.GetCardInfoReq{
		CardKey: d8corp.CardKey{
			Pan:        pan,
			ExpiryDate: expDate,
		},
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b GeneratePIN REQ marshaling err: %v", err)
		return fmt.Errorf("[SERVICE] D8 G2b GeneratePIN REQ marshaling err")
	}
	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/generatePIN", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b GeneratePIN request sending err: %v", err)
		return err
	}
	logger.Infof("[SERVICE] D8 G2b GeneratePIN resp status: %v, body: %v", status, string(data))

	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Errorf("[SERVICE] D8 G2b GeneratePIN RESP marshaling err: %v", err)
		return err
	}
	if resp == nil {
		return fmt.Errorf("[SERVICE] D8 G2b GeneratePIN empty response")
	}
	if resp.Status.Code != "0" {
		logger.Errorf("[SERVICE] D8 G2b GeneratePIN RESP status %s", resp.Status.Code)
		return StatusError(resp.Status)
	}
	return nil
}
