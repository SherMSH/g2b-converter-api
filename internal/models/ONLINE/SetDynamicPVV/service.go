package setdynamicpvv

import (
	"converterapi/internal/config"
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"fmt"
	"strings"
)

// Svc устанавливает PIN карты.
//
// Интерфейс повторяет запрос партнёра, но процессинг D8 динамического PVV не
// поддерживает вовсе: в спецификации 1.80 pvv и pvki доступны только на чтение,
// а из методов работы с PIN есть setPIN, verifyPIN, generatePIN и
// resetCardPINTries. Поэтому запрос выполняется через xmiss/setPIN, и PIN
// меняется постоянно - SingleOperation не действует.
//
// PIN-блок под TPK принять мы не можем: этого ключа у нас нет, расшифровать
// блок нечем, а setPIN требует блок под одноразовым 3DES-ключом, который
// генерируем и заворачиваем в RSA-ключ процессинга мы сами. Поэтому в PINBlock
// ожидается открытый PIN - так же, как его принимает POST /g2b/SetPIN.
func Svc(sb *Body) (soapResp *Envelope, err error) {
	req := sb.SoapRq.Req

	if len(req.PAN) == 0 {
		return nil, fmt.Errorf("wrong mandatory field `fimi1:PAN`")
	}

	pin := strings.TrimSpace(req.PINBlock)
	if pin == "" {
		return nil, fmt.Errorf("PINBlock is empty: сброс динамического PVV процессингом не поддерживается")
	}
	if !isClearPIN(pin) {
		return nil, fmt.Errorf("PINBlock must contain clear PIN (4-12 digits): PIN-блок под рабочим ключом процессинг не принимает")
	}

	if len(req.ExpDate) == 0 {
		req.ExpDate, err = service.GetExpDateByPan(req.PAN)
		if err != nil {
			logger.Warnf("[SERVICE] SetDynamicPVV warning! GetExpDateByPan error: %v", err)
			req.ExpDate = config.Config.Processing.Extra["agreed_expdate_yymm"]
			err = nil
		}
	}

	// О неприменимых полях предупреждаем, но запрос не отклоняем: молча
	// проглотить их хуже - партнёр будет считать, что они сработали
	if req.SingleOperation == "1" {
		logger.Warnf("[SERVICE] SetDynamicPVV: SingleOperation=1 по карте %s не выполнено - процессинг ставит PIN постоянно", req.PAN)
	}
	if req.PVKI != "" && req.PVKI != "0" {
		logger.Warnf("[SERVICE] SetDynamicPVV: PVKI=%s не применён - setPIN его не меняет", req.PVKI)
	}

	if err = service.SetPinG2b(req.PAN, pin, req.ExpDate); err != nil {
		return nil, err
	}

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	soapResp.Body.SetDynamicPVVRp.Response = Response{
		Echo:         req.Echo,
		Product:      req.Product,
		ResponseAttr: "1",
		TranId:       utils.GenerateTimestampID(),
		Ver:          "1.0",
	}
	return soapResp, nil
}

// isClearPIN проверяет, что пришёл именно PIN, а не шифрованный блок.
// Длина PIN по ANSI X9.8 - от 4 до 12 цифр.
func isClearPIN(pin string) bool {
	if len(pin) < 4 || len(pin) > 12 {
		return false
	}
	for _, r := range pin {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
