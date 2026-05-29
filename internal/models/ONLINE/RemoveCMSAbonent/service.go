package removecmsabonent

import (
	models "converterapi/internal/models/OFFLINE"
	service "converterapi/internal/service/G2B"
	"converterapi/pkg/logger"
)

func Svc(sb *Body) (soapResp *Envelope, err error) {
	var address string
	if sb.SoapRq.Req.AlternativeMessaging.Row != nil {
		address = sb.SoapRq.Req.AlternativeMessaging.Row[0].Address
	}
	root := models.Root{
		Records: []models.MRecord{
			{
				PAN:     sb.SoapRq.Req.PAN,
				Address: address,
			},
		},
	}

	result, err := service.DeleteCardNotificationG2b(root)
	logger.Infof("[SERVICE] removecmsabonent result %v", result)
	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	resp := Response{}
	resp.Echo = sb.SoapRq.Req.Echo
	resp.Product = sb.SoapRq.Req.Product
	resp.ResponseAttr = "1"
	resp.Ver = sb.SoapRq.Req.Ver

	soapResp.Body = RespBody{
		RemoveCMSAbonentRp: RemoveCMSAbonentRp{
			Response: resp,
		},
	}
	return soapResp, nil
}
