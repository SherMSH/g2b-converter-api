package setdynamicpvv

import (
	"converterapi/internal/config"
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"fmt"
	"strings"
)

// Svc назначает карте новый PIN.
//
// Интерфейс повторяет запрос партнёра, но процессинг D8 динамического PVV не
// поддерживает вовсе: в спецификации 1.80 pvv и pvki доступны только на чтение.
// Поэтому запрос выполняется через xmiss/generatePIN, и PIN меняется постоянно -
// SingleOperation не действует.
//
// Значение PIN задать нельзя: его выбирает процессинг и наружу не отдаёт
// (спецификация, 7.7). До держателя карты новый PIN доходит SMS-оповещением,
// поэтому у карты должен быть контракт SMSGEN. Присланный PINBlock отклоняем,
// а не игнорируем: иначе на той стороне будут считать, что PIN назначен их.
//
// TODO: ручной ввод PIN вернём. В PINBlock тогда ожидается открытый PIN
// (4-12 цифр) и вызывается xmiss/setPIN вместо generatePIN - реализация была в
// service.SetPinG2b, удалена в 3083cae. PIN-блок под рабочим ключом принять
// по-прежнему нельзя: ключа у нас нет, а setPIN требует блок под одноразовым
// 3DES-ключом, который собираем мы сами.
func Svc(sb *Body) (soapResp *Envelope, err error) {
	req := sb.SoapRq.Req

	if len(req.PAN) == 0 {
		return nil, fmt.Errorf("wrong mandatory field `fimi1:PAN`")
	}

	// TODO: при возврате ручного ввода здесь вместо отказа разбирается PIN
	if strings.TrimSpace(req.PINBlock) != "" {
		return nil, fmt.Errorf("PINBlock must be empty: PIN задаётся процессингом, ручной ввод отключён")
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
		logger.Warnf("[SERVICE] SetDynamicPVV: PVKI=%s не применён - процессинг его не меняет", req.PVKI)
	}

	if err = service.GeneratePIN(req.PAN, req.ExpDate); err != nil {
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
