package getaccinfo

func Svc(sb *Body) (soapResp *Envelope, err error) {
	// accInfo, err := service.GetAccInfo()
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
		GetAcctInfoRp: GetAcctInfoRp{
			Response: resp,
		},
	}
	return soapResp, nil
}
