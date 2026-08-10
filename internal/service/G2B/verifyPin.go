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
	req, err := buildPinRequest(pan, pin, expDate)
	if err != nil {
		return err
	}
	return verifyPinRequest(req)
}

// verifyPinRequest отправляет готовый запрос проверки PIN.
// Вынесено отдельно, чтобы можно было проверить сборку PIN-блока.
func verifyPinRequest(req d8corp.SetPinReq) error {
	var resp *d8corp.CommonResp

	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b verifyPIN REQ marshaling err: %v", err)
		return fmt.Errorf("[SERVICE] D8 G2b verifyPIN REQ marshaling err")
	}

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/verifyPIN", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b verifyPIN request sending err: %v", err)
		return err
	}
	logger.Infof("[SERVICE] D8 G2b verifyPIN resp status: %v, body: %v", status, string(data))

	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Errorf("[SERVICE] D8 G2b verifyPIN RESP marshaling err: %v", err)
		return err
	}

	if resp.Status.Code == "0" && resp.Status.RspCode == "00" {
		return nil
	}
	if reason, ok := pinVerifyErrors[resp.Status.Code+"/"+resp.Status.RspCode]; ok {
		return fmt.Errorf("%s", reason)
	}
	return StatusError(resp.Status)
}
