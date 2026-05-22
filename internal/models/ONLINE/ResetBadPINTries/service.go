package resetbadpintries

import service "converterapi/internal/service/G2B"

func Svc(sb *Body) (soapResp *Envelope, err error) {
	err = service.ResetCardPINTriesG2b(sb.SoapRq.Req.PAN, sb.SoapRq.Req.ExpirationDate)
	if err != nil {
		return nil, err
	}

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	soapResp.Body.ResetBadPINTriesRp.Response = Response{}
	return
}
