package changepin

import (
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"fmt"
)

func Svc(sb *Body) (soapResp *Envelope, err error) {
	//Basic checkup
	if len(sb.SoapRq.Req.PAN) == 0 {
		return nil, fmt.Errorf("400 Bad request. Empty pan")
	}

	expdate, err := service.GetExpDateByPan(sb.SoapRq.Req.PAN)
	service.GeneratePIN(sb.SoapRq.Req.PAN, expdate)

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	resp := Response{}
	resp.Echo = sb.SoapRq.Req.Echo
	resp.Product = sb.SoapRq.Req.Product
	resp.ResponseAttr = "1"
	resp.Ver = sb.SoapRq.Req.Ver
	resp.TranId = utils.GenerateTimestampID()

	soapResp.Body = RespBody{
		ChangePIN: ChangePIN{
			Response: resp,
		},
	}

	return soapResp, nil
}
