package relinkpreissuedcards

import (
	models "converterapi/internal/models/OFFLINE"
	service "converterapi/internal/service/G2B"
)

func Svc(sb *Body) (soapResp *Envelope, err error) {
	root := models.Root{
		Records: []models.MRecord{
			{
				Account:    sb.UpdateCard2AccLink.Req.Account,
				PAN:        sb.UpdateCard2AccLink.Req.PAN,
				MBR:        sb.UpdateCard2AccLink.Req.MBR,
				AcctStat:   sb.UpdateCard2AccLink.Req.AcctStatus,
				Acct2CDesc: sb.UpdateCard2AccLink.Req.Description,
			},
		},
	}

	_, err = service.UpdateCardAcctLinkG2b(root)
	if err != nil {
		return nil, err
	}

	root.Records = []models.MRecord{
		{
			PAN:        sb.UpdateCard2AccLink.Req.PAN,
			MBR:        sb.UpdateCard2AccLink.Req.MBR,
			Account:    sb.UpdateCard2AccLink.Req.Account,
			Acct2CDesc: sb.UpdateCard2AccLink.Req.ChangeReason,
		},
	}
	_, err = service.DeleteCardAcctLinkG2b(root)
	if err != nil {
		return nil, err
	}
	// TODO: Update card for SetCardPerson
	err = service.SetCardStatusG2b(sb.SetCardStatusRq.Req.PAN, sb.SetCardStatusRq.Req.ExpirityDate, sb.SetCardPersonRq.Req.Status, sb.SetCardPersonRq.Req.ChangeReason)
	if err != nil {
		return nil, err
	}

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	resp := Response{
		Echo:         sb.UpdateCard2AccLink.Req.Echo,
		Product:      sb.UpdateCard2AccLink.Req.Product,
		ResponseAttr: "1",
		Ver:          sb.UpdateCard2AccLink.Req.Ver,
		TranId:       "",
	}

	soapResp.Body = RespBody{
		UpdateCard2AcctLinkRp: UpdateCard2AcctLinkRp{
			Response: resp,
		},
		DeleteCard2AcctLinkRp: DeleteCard2AcctLinkRp{
			Response: resp,
		},
		SetCardPersonRp: SetCardPersonRp{
			Response: resp,
		},
		SetCardStatusRp: SetCardStatusRp{
			Response: resp,
		},
	}
	return
}
