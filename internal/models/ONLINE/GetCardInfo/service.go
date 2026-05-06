package getcardinfo

import service "converterapi/internal/service/G2B"

func Svc(sb *Body) (err error) {
	service.GetCardInfo(sb.SoapRq.Req.PAN, sb.SoapRq.Req.ExpirationDate)
	return nil
}
