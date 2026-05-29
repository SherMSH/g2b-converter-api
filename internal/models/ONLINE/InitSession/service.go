package initsession

func Svc(sb *Body) (soapResp *Envelope, err error) {
	// sessionData, err := service.InitiateTransaction()

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	resp := Response{
		Echo:         sb.SoapRq.Req.Echo,
		Product:      sb.SoapRq.Req.Product,
		ResponseAttr: "1",
		Ver:          sb.SoapRq.Req.Ver,

		Id:              "",
		NeedCAPAuth:     "0",
		PasswordVersion: "1",
	}

	resp.CNSChannelList = CNSChannelList{
		Rows: make([]CNSChannelRow, 0),
	}
	resp.CNSChannelList.Rows = append(resp.CNSChannelList.Rows, CNSChannelRow{
		Title:       "ARVAND",
		Name:        "ARVAND",
		UsedForPush: "0",
	})
	resp.CNSSchemeList = CNSSchemeList{
		Rows: make([]CNSSchemeRow, 0),
	}
	resp.CNSSchemeList.Rows = append(resp.CNSSchemeList.Rows, CNSSchemeRow{
		Title: "SMSGEN",
		Name:  "SMSGEN",
	})
	resp.CNSSchemeList.Rows = append(resp.CNSSchemeList.Rows, CNSSchemeRow{
		Title: "SMSTXN",
		Name:  "SMSTXN",
	})

	soapResp.Body = RespBody{
		InitSessionRp: InitSessionRp{
			Response: resp,
		},
	}
	return
}
