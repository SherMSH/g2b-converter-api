package verifypin

import (
	"converterapi/internal/config"
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"fmt"
)

// Svc проверяет PIN по карте.
//
// Неверный PIN - это не ошибка сервиса, а обычный ответ с кодом причины:
// партнёру нужно отличать «PIN не подошёл» от «проверить не удалось», поэтому
// SOAP Fault остаётся только для некорректного запроса.
//
// Внимание: неудачные попытки увеличивают счётчик неверных вводов в процессинге
// и в итоге блокируют карту. Сброс - операцией ResetBadPINTriesRq.
func Svc(sb *Body) (soapResp *Envelope, err error) {
	req := sb.SoapRq.Req

	if len(req.PAN) == 0 {
		return nil, fmt.Errorf("Bad Request: empty PAN field")
	}
	if len(req.PIN) == 0 {
		return nil, fmt.Errorf("Bad Request: empty PIN field")
	}
	if len(req.ExpirationDate) == 0 {
		req.ExpirationDate, err = service.GetExpDateByPan(req.PAN)
		if err != nil {
			logger.Warnf("[SERVICE] VerifyPIN warning! GetExpDateByPan error: %v", err)
			req.ExpirationDate = config.Config.Processing.Extra["agreed_expdate_yymm"]
			err = nil
		}
	}

	status, err := service.VerifyPinStatusG2b(req.PAN, req.PIN, req.ExpirationDate)
	if err != nil {
		return nil, err
	}

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	soapResp.Body.VerifyPINRp.Response = Response{
		Echo:         req.Echo,
		Product:      req.Product,
		ResponseAttr: "1",
		TranId:       utils.GenerateTimestampID(),
		Ver:          "1.0",
		AuthRespCode: utils.AuthRespCode(status.Code, status.RspCode),
	}
	if !utils.IsApproved(status.Code) {
		soapResp.Body.VerifyPINRp.Response.DeclineReason = service.PinVerifyReason(status)
	}
	return soapResp, nil
}
