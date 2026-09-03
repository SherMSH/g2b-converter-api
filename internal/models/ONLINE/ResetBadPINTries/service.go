package resetbadpintries

import (
	"converterapi/internal/config"
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"fmt"
)

func Svc(sb *Body) (soapResp *Envelope, err error) {
	if len(sb.SoapRq.Req.PAN) == 0 {
		return nil, fmt.Errorf("Bad Request: empty PAN field")
	}
	if len(sb.SoapRq.Req.ExpirationDate) == 0 {
		sb.SoapRq.Req.ExpirationDate, err = service.GetExpDateByPan(sb.SoapRq.Req.PAN)
		if err != nil {
			logger.Warnf("[SERVICE] GetCardInfo warning! GetExpDateByPan error: %v", err)
			sb.SoapRq.Req.ExpirationDate = config.Config.Processing.Extra["agreed_expdate_yymm"]
			err = nil
		}
	}
	err = service.ResetCardPINTriesG2b(sb.SoapRq.Req.PAN, sb.SoapRq.Req.ExpirationDate)
	if err != nil {
		return nil, err
	}

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	soapResp.Body.ResetBadPINTriesRp.Response = Response{
		Echo:         sb.SoapRq.Req.Echo,
		Product:      sb.SoapRq.Req.Product,
		ResponseAttr: "1",
		TranId:       utils.GenerateTimestampID(),
		Ver:          "1.0",
	}
	return
}
