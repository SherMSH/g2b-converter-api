package getcardstatement

import (
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"fmt"
	"strconv"
	"time"
)

func Svc(sb *Body) (soapResp *Envelope, err error) {
	from, to := getFormatedTime(sb.SoapRq.Req.FromTime, sb.SoapRq.Req.ToTime)
	size := 10
	if len(sb.SoapRq.Req.Count) > 0 {
		size, err = strconv.Atoi(sb.SoapRq.Req.Count)
		if err != nil {
			logger.Errorf("[SERVICE] getcardstatement req error: wrong Count param")
			err = nil
		}
	}
	cardInfo, err := service.GetCardTransactionHistory(sb.SoapRq.Req.PAN, from, to, size)
	if err != nil {
		return nil, err
	}

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	resp := Response{}
	resp.Product = sb.SoapRq.Req.Product
	resp.Echo = sb.SoapRq.Req.Echo
	resp.ResponseAttr = "1"
	resp.Ver = sb.SoapRq.Req.Ver
	resp.TranId = utils.GenerateTimestampID()

	rows := make([]Row, 0)
	for i, v := range cardInfo.CardTransactions {
		operDate, _ := time.ParseInLocation("20060102", v.BusDate, time.Local)
		tranTime, _ := time.ParseInLocation("20060102150405", v.When_created[:14], time.Local)
		row := Row{
			FrontId:             fmt.Sprintf("%v", v.TlId),
			Origin:              fmt.Sprintf("%d", i),
			Type:                "1",
			OperCode:            fmt.Sprintf("%d", v.Txncode),
			Description:         "",
			Amount:              fmt.Sprintf("%.2f", v.TxnAmount),
			Currency:            utils.Currencies[v.TxnCurrency],
			OperDate:            operDate.Format("2006-01-02"),
			TranTime:            tranTime.Format("2006-01-02T15:04:05"),
			OrigAmount:          fmt.Sprintf("%.2f", v.Amtbill),
			OrigCurrency:        utils.Currencies[v.Curbill],
			PAN:                 v.Lkey.Pan,
			MBR:                 "0",
			TermClass:           v.TermType,
			TermName:            fmt.Sprintf("%s %s", v.CrdacptID, v.TermCode),
			TermRetailerName:    v.CrdactplocName,
			TermSIC:             fmt.Sprintf("%d", v.CrdacptBus),
			TermLocation:        v.CrdactplocName,
			ApprovalCode:        fmt.Sprintf("%.6d", v.Stan/1000000),
			CurrencyISOCode:     v.TxnCurrency,
			OrigCurrencyISOCode: v.Curbill,
			SeqNo:               fmt.Sprintf("%v", i),
			OrigTime:            tranTime.Format("2006-01-02T15:04:05"),
			TermCountry:         v.CrdactplocCountry,
			TermCity:            v.CrdactplocCity,
			// OnlineIssuerFee:     "",
		}
		rows = append(rows, row)
	}

	resp.Statement = Statement{
		Rows: rows,
	}

	soapResp.Body = RespBody{
		GetCardStatementRp: GetCardStatementRp{
			Response: resp,
		},
	}
	return soapResp, nil
}

func getFormatedTime(reqFrom, reqTo string) (from, to string) {
	tFrom, err := time.ParseInLocation("2006-01-02T15:04:05", reqFrom, time.Local)
	if err != nil {
		return "", ""
	}
	tTo, err := time.ParseInLocation("2006-01-02T15:04:05", reqTo, time.Local)
	if err != nil {
		return "", ""
	}

	from = tFrom.Format("20060102150405")
	to = tTo.Format("20060102150405")
	return
}
