package service

import (
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/crypto"
	"converterapi/pkg/logger"
	"crypto/rsa"
	"encoding/json"
	"fmt"
)

func SetPinG2b(pan, pin, expDate string) error {
	var resp *d8corp.CommonResp

	pinBlockBuilder := crypto.NewPinBlockBuilder()
	pinblock0, err := pinBlockBuilder.BuildFormat0(pan, pin)
	if err != nil {
		return err
	}

	pinBlockEncrypter := crypto.NewTripleDESEncrypter()
	pinBlock, err := pinBlockEncrypter.EncryptECB(crypto.RandomZPK, pinblock0)
	if err != nil {
		return err
	}

	pubKey := new(rsa.PublicKey)
	pinBlockRSAEncrypter := crypto.NewRSAZPKEncrypter(pubKey)

	pinBlockUnderRSA, err := pinBlockRSAEncrypter.Encrypt(pinBlock)
	if err != nil {
		return err
	}
	req := d8corp.SetPinReq{
		CardKey: d8corp.CardKey{
			Pan:        pan,
			ExpiryDate: expDate,
		},
		PinKeyUnderRSA: string(pinBlockUnderRSA),
		PinBlock:       string(pinBlock),
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
		return fmt.Errorf("%s - %s", resp.Status.RspCode, resp.Status.Message)
	}
	// err = json.Unmarshal(resp.Data, &respData)
	// if err != nil {
	// 	logger.Errorf("[SERVICE] D8 G2b setPIN RESP data marshaling err: %v", err)
	// 	return err
	// }
	// if respData == nil {
	// 	logger.Errorf("[SERVICE] D8 G2b setPIN RESP data is empty")
	// 	return fmt.Errorf("no data")
	// }
	return nil
}

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
