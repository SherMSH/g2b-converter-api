package resetbadpintries

import (
	"converterapi/internal/config"
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
)

func Svc(sb *Body) (soapResp *Envelope, err error) {
	if config.Config.App.DebugMode && len(sb.SoapRq.Req.ExpirationDate) == 0 {
		sb.SoapRq.Req.ExpirationDate = "3004"
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
