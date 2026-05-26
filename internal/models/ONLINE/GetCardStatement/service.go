package getcardstatement

import (
	service "converterapi/internal/service/G2B"
	"fmt"
	"time"
)

func Svc(sb *Body) (soapResp *Envelope, err error) {
	cardInfo, err := service.GetCardInfo(sb.SoapRq.Req.PAN, sb.SoapRq.Req.ExpirationDate)
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

	rows := make([]Row, 0)
	for i, v := range cardInfo.CardTransactions {
		operDate, _ := time.ParseInLocation("20060102", v.BusDate, time.Local)
		tranTime, _ := time.ParseInLocation("20060102150405", v.When_created, time.Local)
		row := Row{
			FrontId:             fmt.Sprintf("%v", v.TlId),
			Origin:              fmt.Sprintf("%d", i),
			Type:                "1",
			OperCode:            fmt.Sprintf("%d", v.Txncode),
			Description:         "",
			Amount:              fmt.Sprintf("%.2f", v.TxnAmount),
			Currency:            v.TxnCurrency,
			OperDate:            operDate.Format("2006-01-02"),
			TranTime:            tranTime.Format("2006-01-02T15:04:05"),
			OrigAmount:          fmt.Sprintf("%.2f", v.Amtbill),
			OrigCurrency:        v.Curbill,
			PAN:                 v.Lkey.MaskedPan,
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
