package acctcredit

import "converterapi/internal/utils"

func Svc(sb *Body) (soapResp *Envelope, err error) {
	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	resp := Response{}
	resp.Product = sb.SoapRq.Req.Product
	resp.ResponseAttr = "1"
	resp.Ver = sb.SoapRq.Req.Ver
	resp.TranId = utils.GenerateTimestampID()

	soapResp.Body = RespBody{
		AcctCreditRp: AcctCreditRp{
			Response: resp,
		},
	}
	return soapResp, nil
}
