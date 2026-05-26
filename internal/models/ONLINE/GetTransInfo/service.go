package gettransinfo

import service "converterapi/internal/service/G2B"

func Svc(b *Body) (soapResp *Envelope, err error) {
	_ = service.GetTransactionInfoG2b()
	if err != nil {
		return nil, err
	}

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	tranListArr := []TranListRow{
		TranListRow{
			TranNumber: b.SoapRq.Req.TranNumber,
			PAN:        b.SoapRq.Req.PAN,
			TermName:   b.SoapRq.Req.TermName,
			TranCode:   b.SoapRq.Req.TranCode.Row[0].Value,
		},
	}

	soapResp.Body.GetTransInfoRp.Response = Response{
		Echo:         b.SoapRq.Req.Echo,
		Product:      b.SoapRq.Req.Product,
		ResponseAttr: "1",
		Ver:          "1.0",

		MaskBalances: "0",
		TranList: TranList{
			Rows: tranListArr,
		},
	}
	return
}
