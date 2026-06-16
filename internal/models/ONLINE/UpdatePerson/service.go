package updateperson

import (
	models "converterapi/internal/models/OFFLINE"
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
)

func Svc(sb *Body) (soapResp *Envelope, err error) {

	root := models.Root{
		Records: []models.MRecord{
			{
				Id:         sb.SoapRq.Req.Id,
				FirstName:  sb.SoapRq.Req.FirstName,
				LastName:   sb.SoapRq.Req.LastName,
				MiddleName: sb.SoapRq.Req.MiddleName,
				Birthday:   sb.SoapRq.Req.Birthday,
				CntrLive:   sb.SoapRq.Req.ResidentCountry,
			},
		},
	}

	_, err = service.UpdateCustomerG2b(root)
	if err != nil {
		return nil, err
	}
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
		UpdatePersonRp: UpdatePersonRp{
			Response: resp,
		},
	}
	return soapResp, nil
}
