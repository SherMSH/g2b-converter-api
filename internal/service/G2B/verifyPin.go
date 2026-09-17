package service

import (
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
)

// pinVerifyErrors - коды ответа xmiss/verifyPIN (7.22.2).
// Ключ - "code/rspcode".
var pinVerifyErrors = map[string]string{
	"1/00": "PIN received when no PIN required",
	"1/06": "PIN tries exceeds max number of allowed tries",
	"1/12": "PIN is missing",
	"1/17": "bad PIN",
	"1/18": "card configuration issue",
	"1/26": "invalid PIN block",
	"9/09": "system error",
}

// VerifyPinG2b проверяет PIN по карте (xmiss/verifyPIN, 7.22).
//
// Успех - status 0/00. Остальные коды разложены в понятные сообщения:
// по ним видно, отличается ли неверный PIN от испорченного PIN-блока.
//
// Внимание: неудачные попытки увеличивают счётчик неверных вводов и в итоге
// блокируют карту - сбрасывается через ResetCardPINTriesG2b.
func VerifyPinG2b(pan, pin, expDate string) error {
	status, err := VerifyPinStatusG2b(pan, pin, expDate)
	if err != nil {
		return err
	}
	if status.Code == "0" && status.RspCode == "00" {
		return nil
	}
	return fmt.Errorf("%s", PinVerifyReason(status))
}

// VerifyPinStatusG2b проверяет PIN и отдаёт ответ процессинга как есть.
//
// Нужен операции VerifyPINRq: партнёру уходит не «да/нет», а код причины,
// поэтому схлопывать ответ в ошибку здесь нельзя.
func VerifyPinStatusG2b(pan, pin, expDate string) (d8corp.Status, error) {
	req, err := buildPinRequest(pan, pin, expDate)
	if err != nil {
		return d8corp.Status{}, err
	}
	return verifyPinRequest(req)
}

// PinVerifyReason - человекочитаемая причина отказа.
// Процессинг не всегда присылает message, поэтому подстраховываемся
// собственным справочником кодов.
func PinVerifyReason(status d8corp.Status) string {
	if reason, ok := pinVerifyErrors[status.Code+"/"+status.RspCode]; ok {
		return reason
	}
	if status.Message != "" {
		return status.Message
	}
	return StatusError(status).Error()
}

// verifyPinRequest отправляет готовый запрос проверки PIN.
// Вынесено отдельно, чтобы можно было проверить сборку PIN-блока.
func verifyPinRequest(req d8corp.SetPinReq) (d8corp.Status, error) {
	var resp *d8corp.CommonResp

	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b verifyPIN REQ marshaling err: %v", err)
		return d8corp.Status{}, fmt.Errorf("[SERVICE] D8 G2b verifyPIN REQ marshaling err")
	}

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/verifyPIN", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b verifyPIN request sending err: %v", err)
		return d8corp.Status{}, err
	}
	logger.Infof("[SERVICE] D8 G2b verifyPIN resp status: %v, body: %v", status, string(data))

	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Errorf("[SERVICE] D8 G2b verifyPIN RESP marshaling err: %v", err)
		return d8corp.Status{}, err
	}
	if resp == nil {
		return d8corp.Status{}, fmt.Errorf("[SERVICE] D8 G2b verifyPIN empty response")
	}
	return resp.Status, nil
}
