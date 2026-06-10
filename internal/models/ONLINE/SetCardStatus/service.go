package setcardstatus

import (
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
)

func Svc(sb *Body) (soapResp *Envelope, err error) {
	err = service.SetCardStatusG2b(sb.SoapRq.Req.PAN, sb.SoapRq.Req.ExpirationDate, sb.SoapRq.Req.Status, sb.SoapRq.Req.ChangeReason)
	if err != nil {
		return nil, err
	}

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	soapResp.Body.SetCardStatusRp.Response = Response{
		Echo:         sb.SoapRq.Req.Echo,
		Product:      sb.SoapRq.Req.Product,
		ResponseAttr: "1",
		TranId:       utils.GenerateTimestampID(),
		Ver:          "1.0",
	}
	return
}
