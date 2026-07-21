package changecmsabonent

import (
	"converterapi/internal/config"
	models "converterapi/internal/models/OFFLINE"
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"fmt"
)

func Svc(b *Body) (soapResp *Envelope, err error) {
	if len(b.SoapRq.Req.AltMessaging.Row) == 0 || len(b.SoapRq.Req.AltMessaging.Row[0].Address) == 0 {
		return nil, fmt.Errorf("Mandatory field is empty 'AlternativeMassaging -> Row -> Address'")
	}
	if len(b.SoapRq.Req.AltMessaging.Row[0].PrevAddress) == 0 {
		return nil, fmt.Errorf("Mandatory field is empty 'AlternativeMassaging -> Row -> PrevAddress'")
	}
	if config.Config.App.DebugMode && len(b.SoapRq.Req.ExpirationDate) == 0 {
		b.SoapRq.Req.ExpirationDate = "20300430"
	}

	root := models.Root{}
	record := models.MRecord{
		PAN:     b.SoapRq.Req.PAN,
		ExpDate: b.SoapRq.Req.ExpirationDate,
		Address: b.SoapRq.Req.AltMessaging.Row[0].Address,
	}
	root.Records = append(root.Records, record)

	_, err = service.DeleteCardNotificationG2b(root)
	if err != nil {
		logger.Errorf("%s", err.Error())
	}

	_, err = service.AddCardNotificationG2b(root)
	if err != nil {
		logger.Errorf("%s", err.Error())
		return nil, err
	}

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	soapResp.Body.ChangeCMSAbonentRp.Response = Response{
		Echo:         b.SoapRq.Req.Echo,
		Product:      b.SoapRq.Req.Product,
		ResponseAttr: "1",
		TranId:       utils.GenerateTimestampID(),
		Ver:          "1.0",
	}
	return
}
