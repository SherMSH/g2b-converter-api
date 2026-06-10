package getacctstatement

import (
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
)

func Svc(sb *Body) (soapResp *Envelope, err error) {
	err = service.GetAcctInfoG2b()
	if err != nil {
		return nil, err
	}

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	resp := Response{
		Echo:         sb.SoapRq.Req.Echo,
		Product:      sb.SoapRq.Req.Product,
		ResponseAttr: "1",
		TranId:       utils.GenerateTimestampID(),
		Ver:          "1.0",
	}
	rows := make([]Row, 0)
	rows = append(rows, Row{})

	resp.Statement = Statement{
		Rows: rows,
	}

	soapResp.Body = RespBody{
		GetAcctStatementRp: GetAcctStatementRp{
			Response: resp,
		},
	}
	return soapResp, nil
}
