package setacctstatus

import (
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"fmt"
)

// Svc меняет статус карточного счёта.
//
// Партнёр присылает статус в кодировке TWO, процессинг ждёт свой код, поэтому
// переводим через справочник.
func Svc(sb *Body) (soapResp *Envelope, err error) {
	req := sb.SoapRq.Req

	if len(req.Account) == 0 {
		return nil, fmt.Errorf("wrong mandatory field `fimi1:Account`")
	}
	statCode, ok := utils.ReverseAccountStatuses[req.Status]
	if !ok {
		return nil, fmt.Errorf("Status %v is not supported!", req.Status)
	}

	if err = service.SetAccountStatusG2b(req.Account, statCode); err != nil {
		return nil, err
	}

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	soapResp.Body.SetAcctStatusRp.Response = Response{
		Echo:         req.Echo,
		Product:      req.Product,
		ResponseAttr: "1",
		TranId:       utils.GenerateTimestampID(),
		Ver:          "1.0",
	}
	return soapResp, nil
}
