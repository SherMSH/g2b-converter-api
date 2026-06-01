package getacctstatement

import (
	service "converterapi/internal/service/G2B"
	"fmt"
	"strconv"
	"time"
)

func Svc(sb *Body) (soapResp *Envelope, err error) {
	err = service.GetAcctInfoG2b()
	if err != nil {
		return nil, err
	}
	dateFrom, _ := time.ParseInLocation("2006-01-02T15:04:05", sb.SoapRq.Req.FromTime, time.Local)
	dateTo, _ := time.ParseInLocation("2006-01-02T15:04:05", sb.SoapRq.Req.ToTime, time.Local)
	trnTime := time.Date(2026, 5, 26, 9, 37, 54, 0, time.Local)

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	resp := Response{
		Echo:         sb.SoapRq.Req.Echo,
		Product:      sb.SoapRq.Req.Product,
		ResponseAttr: "1",
		TranId:       "",
		Ver:          "1.0",
	}
	rows := make([]Row, 0)

	rows = append(rows, Row{
		FrontId:      "167",
		Origin:       "XAPI/00005GSLZE4o2Ddn2iCG7KrtMmnxg5Va",
		Type:         "0",
		OperCode:     "00",
		Description:  "*Dushanbe RRP XAPI/00005GSLZE4o2Ddn2iCG7KrtMmnxg5Va",
		Amount:       "1.0",
		OperDate:     "2026-05-26T00:00:00",
		TranTime:     "2026-05-26T09:37:54",
		OrigAmount:   "1.0",
		OrigCurrency: "972",
		PAN:          "5058270310000020",
		MBR:          "0",
		TermClass:    "M",
		TermName:     "Test001",
		TermSIC:      "5999",
		TermLocation: "*Dushanbe RRP",
		// ApprovalCode:        "",
		// SeqNo:               "",
		TermCountry:         "762",
		TermCity:            "Dushanbe",
		OnlineIssuerFee:     "0",
		OrigTime:            "2026-05-26T14:37:54",
		Currency:            "972",
		TermRetailerName:    "",
		CurrencyISOCode:     "972",
		OrigCurrencyISOCode: "972",
	})

	k, _ := strconv.Atoi(sb.SoapRq.Req.Count)
	for i := 0; i < k; i++ {
		if trnTime.After(dateFrom) && trnTime.Before(dateTo) {
			resp.Statement.Rows = append(resp.Statement.Rows, rows[0])
			resp.Statement.Rows[i].SeqNo = fmt.Sprintf("%d", i+1)
		}
	}

	soapResp.Body = RespBody{
		GetAcctStatementRp: GetAcctStatementRp{
			Response: resp,
		},
	}
	return soapResp, nil
}
