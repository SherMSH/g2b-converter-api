package addcmsabonent

import (
	models "converterapi/internal/models/OFFLINE"
	service "converterapi/internal/service/G2B"
	"converterapi/pkg/logger"
	"fmt"
)

func Svc(b *Body) (soapResp *Envelope, err error) {
	if len(b.SoapRq.Req.AltMessaging.Row) == 0 || len(b.SoapRq.Req.AltMessaging.Row[0].Address) == 0 {
		return nil, fmt.Errorf("Mandatory field is empty 'AlternativeMassaging -> Row -> Address'")
	}

	root := models.Root{}
	record := models.MRecord{
		PAN:     b.SoapRq.Req.PAN,
		Address: b.SoapRq.Req.AltMessaging.Row[0].Address,
	}
	root.Records = append(root.Records, record)

	_, err = service.AddCardNotificationG2b(root)
	if err != nil {
		logger.Errorf("%s", err.Error())
		return nil, err
	}

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	soapResp.Body.AddCMSAbonentRp.Response = Response{
		Echo:         b.SoapRq.Req.Echo,
		Product:      b.SoapRq.Req.Product,
		ResponseAttr: "1",
		TranId:       "",
		Ver:          "1.0",
	}

	return
}
