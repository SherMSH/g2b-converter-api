package changepin

import (
	"converterapi/internal/config"
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"fmt"
)

// Svc назначает карте новый PIN.
//
// Значение выбирает процессинг и наружу не отдаёт - до держателя новый PIN
// доходит SMS-оповещением, поэтому у карты должен быть контракт SMSGEN.
func Svc(sb *Body) (soapResp *Envelope, err error) {
	//Basic checkup
	if len(sb.SoapRq.Req.PAN) == 0 {
		return nil, fmt.Errorf("400 Bad request. Empty pan")
	}

	expdate, err := service.GetExpDateByPan(sb.SoapRq.Req.PAN)
	if err != nil {
		logger.Warnf("[SERVICE] ChangePIN warning! GetExpDateByPan error: %v", err)
		expdate = config.Config.Processing.Extra["agreed_expdate_yymm"]
	}

	// Ошибку генерации нельзя терять: иначе партнёр получит успех, а PIN
	// останется прежним
	if err = service.GeneratePIN(sb.SoapRq.Req.PAN, expdate); err != nil {
		return nil, err
	}

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	resp := Response{}
	resp.Echo = sb.SoapRq.Req.Echo
	resp.Product = sb.SoapRq.Req.Product
	resp.ResponseAttr = "1"
	resp.Ver = sb.SoapRq.Req.Ver
	resp.TranId = utils.GenerateTimestampID()

	soapResp.Body = RespBody{
		ChangePIN: ChangePIN{
			Response: resp,
		},
	}

	return soapResp, nil
}
