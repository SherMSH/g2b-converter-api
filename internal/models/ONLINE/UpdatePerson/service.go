package updateperson

func Svc(sb *Body) (soapResp *Envelope, err error) {
	// cardInfo, err := service.GetCardInfo(sb.SoapRq.Req.PAN, sb.SoapRq.Req.ExpirationDate)
	// if err != nil {
	// 	return nil, err
	// }
	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	resp := Response{}
	resp.Product = sb.SoapRq.Req.Product
	resp.ResponseAttr = "1"
	resp.Ver = sb.SoapRq.Req.Ver

	soapResp.Body = RespBody{
		UpdatePersonRp: UpdatePersonRp{
			Response: resp,
		},
	}
	return soapResp, nil
}
