package getpersoninfo

import (
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
)

func Svc(b *Body) (soapResp *Envelope, err error) {
	err = service.GetCustomerInfoG2b()
	if err != nil {
		return nil, err
	}
	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	soapResp.Body.GetPersonInfoRp.Response = Response{
		Echo:         b.SoapRq.Req.Echo,
		Product:      b.SoapRq.Req.Product,
		ResponseAttr: "1",
		TranId:       utils.GenerateTimestampID(),
		Ver:          "1.0",
	}

	soapResp.Body.GetPersonInfoRp.Response.Accounts = []Account{
		{
			PersonId: b.SoapRq.Req.Id,
			Account:  "20216972600001090524",
			Type:     "0",
			Status:   "00 - Active",
		},
	}

	soapResp.Body.GetPersonInfoRp.Response.AlternativeMessaging = []AlternativeMessage{
		{
			Channel:        "SMSGEN",
			Address:        "006772212",
			Scheme:         "ARVD",
			Disabled:       "0",
			UseForDynAuth:  "1",
			IsDefault:      "1",
			Broadcast:      "1",
			DynAuthBlocked: "0",
			ErrorCount:     "0",
			CreationDate:   "21-04-2026T15:04:05",
			LastUpdateTime: "21-04-2026T15:04:05",
		},
	}

	soapResp.Body.GetPersonInfoRp.Response.Cards = []Card{
		{
			PersonId: b.SoapRq.Req.Id,
			PAN:      "5058270530000032",
			MBR:      "0",
			Status:   "00",
			ExpDate:  "2026-05-31T23:59:59",
			Type:     "0",
		},
	}

	soapResp.Body.GetPersonInfoRp.Response.Confidential = []Confidential{
		{
			PersonId:     b.SoapRq.Req.Id,
			What:         "word",
			Value:        "confidential_word",
			IsAllowedCST: "0",
			IsAllowedADS: "0",
			IsAllowedTB:  "0",
		},
	}

	soapResp.Body.GetPersonInfoRp.Response.Info = []Info{
		{
			PersonId:        b.SoapRq.Req.Id,
			FIO:             "006772212 Давронбек Болтабоев",
			Sex:             "M",
			IdentType:       "1",
			Identity:        "A21404440",
			Birthday:        "1991-07-19T00:00:00",
			Birthplace:      "Hudjand",
			VIP:             "1",
			InstName:        "ARVD",
			PersonExtId:     "120147",
			ResidentCountry: "972",
			AddressInLatin:  "",
			FirstNameNat:    "Давронбек",
			LastNameNat:     "006772212",
			MiddleNameNat:   "Болтабоев",
			TaxPayerNumber:  "",
			AddressNat:      "",
		},
	}

	return
}
