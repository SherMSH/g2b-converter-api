package acctdebit

func Svc(b *Body) (soapResp *Envelope, err error) {
	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	soapResp.Body.AcctDebitRp.Response = Response{
		Product:      b.SoapRq.Req.Product,
		ResponseAttr: "1",
		TranId:       "",
		Ver:          "1.0",

		ApprovalCode:  "780D1O",
		AvailBalance:  321.22,
		LedgerBalance: 0.00,
	}
	return
}
