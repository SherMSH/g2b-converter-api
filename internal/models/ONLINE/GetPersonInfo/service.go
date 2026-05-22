package getpersoninfo

import service "converterapi/internal/service/G2B"

func Svc(b *Body) (soapResp *Envelope, err error) {
	err = service.GetCustomerInfoG2b()
	if err != nil {
		return nil, err
	}
	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	soapResp.Body.GetPersonInfoRp.Response = Response{
		Echo:         b.SoapRq.Req.Echo,
		Product:      b.SoapRq.Req.Product,
		ResponseAttr: "1",
		TranId:       "",
		Ver:          "1.0",
	}
	return
}
