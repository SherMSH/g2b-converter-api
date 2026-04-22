package changecmsabonent

import (
	models "converterapi/internal/models/OFFLINE"
	service "converterapi/internal/service/G2B"
	"converterapi/pkg/logger"
	"fmt"
)

func Svc(b *Body) error {
	if len(b.SoapRq.Req.AltMessaging.Row) == 0 || len(b.SoapRq.Req.AltMessaging.Row[0].Address) == 0 {
		return fmt.Errorf("Mandatory field is empty 'AlternativeMassaging -> Row -> Address'")
	}
	if len(b.SoapRq.Req.AltMessaging.Row[0].PrevAddress) == 0 {
		return fmt.Errorf("Mandatory field is empty 'AlternativeMassaging -> Row -> PrevAddress'")
	}

	root := models.Root{}
	record := models.MRecord{
		PAN:     b.SoapRq.Req.PAN,
		Address: b.SoapRq.Req.AltMessaging.Row[0].Address,
	}
	root.Records = append(root.Records, record)

	_, err := service.AddCardNotificationG2b(root)
	if err != nil {
		logger.Errorf("%s", err.Error())
		return err
	}

	_, err = service.DeleteCardNotificationG2b(root)
	if err != nil {
		logger.Errorf("%s", err.Error())
	}

	return nil
}
