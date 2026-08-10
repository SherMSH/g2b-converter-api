package service

import (
	"converterapi/internal/config"
	"converterapi/internal/models"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
)

// CalculateAcqCommission рассчитывает комиссию эквайрера по операции
// (xkernel/calculateAcqCommission, 8.8).
//
// Запрос идентичен авторизации, поэтому собирается тем же BuildAuthTxReq.
// Как и авторизации, расчёту нужна ссылка от InitiateTransaction.
//
// Средства не списываются: это только расчёт.
func CalculateAcqCommission(input models.TrnInputIface, ecTxRefNo string) (commission *d8corp.AcqCommissionData, err error) {
	var resp *d8corp.CommonResp

	req, err := BuildAuthTxReq(input, ecTxRefNo)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b calculateAcqCommission req err: %v", err)
		return nil, err
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b calculateAcqCommission REQ marshaling err: %v", err)
		return nil, fmt.Errorf("[SERVICE] D8 G2b calculateAcqCommission REQ marshaling err")
	}

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/kernel/1.0/calculateAcqCommission", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b calculateAcqCommission request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b calculateAcqCommission resp status: %v, body: %v (req %v)", status, string(data), string(jsonReq))

	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Errorf("[SERVICE] D8 G2b calculateAcqCommission RESP marshaling err: %v", err)
		return nil, err
	}
	if resp.Status.Code != "0" {
		logger.Errorf("[SERVICE] D8 G2b calculateAcqCommission RESP status %+v", resp.Status)
		return nil, StatusError(resp.Status)
	}
	if err = json.Unmarshal(resp.Data, &commission); err != nil {
		logger.Errorf("[SERVICE] D8 G2b calculateAcqCommission DATA marshaling err: %v", err)
		return nil, err
	}
	if commission == nil {
		return nil, fmt.Errorf("no data")
	}
	return commission, nil
}

// StatusError превращает статус ответа D8 в ошибку.
//
// Процессинг не всегда заполняет message: у отказа расчёта комиссии по
// переводу приходит только code=1, и привычное "rspcode - message" даёт
// бесполезное "00 - ". В таком случае показываем сам код.
func StatusError(status d8corp.Status) error {
	if status.Message != "" {
		return fmt.Errorf("%s - %s", status.RspCode, status.Message)
	}
	return fmt.Errorf("processing status code %s (rspcode %s)", status.Code, status.RspCode)
}
