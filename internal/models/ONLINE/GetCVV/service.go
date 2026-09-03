package getcvv

import (
	"converterapi/internal/config"
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"fmt"
)

func Svc(sb *Body) (soapResp *Envelope, err error) {
	if len(sb.SoapRq.Req.PAN) == 0 {
		return nil, fmt.Errorf("Bad Request: empty PAN field")
	}
	if len(sb.SoapRq.Req.ExpDate) == 0 {
		sb.SoapRq.Req.ExpDate, err = service.GetExpDateByPan(sb.SoapRq.Req.PAN)
		if err != nil {
			logger.Warnf("[SERVICE] GetCVV warning! GetExpDateByPan error: %v", err)
			sb.SoapRq.Req.ExpDate = config.Config.Processing.Extra["agreed_expdate_yymm"]
			err = nil
		}
	}
	cvvData, err := service.GetCVVG2b(sb.SoapRq.Req.PAN, sb.SoapRq.Req.ExpDate)
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

	resp.CVV = cvvData.CVV2
	resp.CardVerificationType = fmt.Sprintf("%d", cvvData.CVV2Type)
	resp.StrCVV = cvvData.CVV2

	soapResp.Body.GetCVVRp.Response = resp
	return
}
