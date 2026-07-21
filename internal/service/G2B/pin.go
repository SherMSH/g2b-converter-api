package service

import (
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/crypto"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
)

func SetPinG2b(pan, pin, expDate string) error {
	var resp *d8corp.CommonResp

	key3DES, err := crypto.Generate3DESKey()
	if err != nil {
		logger.Errorf("generate 3DES key error: %v", err)
		return fmt.Errorf("generate 3DES key error: %v", err)
	}
	logger.Infof("Generated 3DES key: %v", crypto.HexUpper(key3DES))

	publicKey, err := crypto.ReadPublicKey("internal/app/files/transport_setpin.der")
	if err != nil {
		logger.Errorf("read public key error: %v", err)
		return fmt.Errorf("read public key error: %v", err)
	}
	logger.Infof("Key size: %d", publicKey.N.BitLen())

	// Шифруем 3DES ключ публичным ключом
	pinKeyUnderRSA, err := crypto.EncryptWithRSA(publicKey, key3DES)
	if err != nil {
		logger.Errorf("Encrypt with RSA: %v", err)
		return fmt.Errorf("Encrypt with RSA: %v", err)
	}

	clear, err := crypto.Format0(pin, pan)
	if err != nil {
		logger.Errorf("pin block format0: %v", err)
		return fmt.Errorf("pin block format0: %v", err)
	}
	var zpk, encrypted []byte
	zpk, err = crypto.GenerateZPK32()
	if err != nil {
		logger.Errorf("generate ZPK: %v", err)
		return fmt.Errorf("generate ZPK: %v", err)
	}
	encrypted, err = crypto.Encrypt3DES(zpk, clear)
	if err != nil {
		logger.Errorf("encrypt: %v", err)
		return fmt.Errorf("encrypt: %v", err)
	}

	pinBlock := encrypted // hex.DecodeString(pinBlockHex)
	logger.Infof("pinBlock: %X", pinBlock)

	// Шифруем PIN-блок 3DES ключом
	encryptedPinBlock, err := crypto.EncryptWith3DES(key3DES, pinBlock)
	if err != nil {
		logger.Errorf("Encrypt with 3DES: %v", err)
		return fmt.Errorf("Encrypt with 3DES: %v", err)
	}

	req := d8corp.SetPinReq{
		CardKey: d8corp.CardKey{
			Pan:        pan,
			ExpiryDate: expDate,
		},
		PinKeyUnderRSA: crypto.HexUpper(pinKeyUnderRSA),
		PinBlock:       crypto.HexUpper(encryptedPinBlock),
		PinBlockType:   0,
	}

	logger.Infof("[SERVICE] setPin request: %+v", req)

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
