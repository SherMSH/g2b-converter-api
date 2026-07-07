package getacctstatement

import (
	d8corp "converterapi/internal/models/D8CORP"
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"fmt"
	"strconv"
	"time"
)

func Svc(sb *Body) (soapResp *Envelope, err error) {
	dateFrom, errPrsFrom := time.ParseInLocation("2006-01-02T15:04:05", sb.SoapRq.Req.FromTime, time.Local)
	if errPrsFrom != nil {
		dateFrom = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		logger.Errorf("Date `From` parsing error: %v; setting default: %v", err, dateFrom.Format("2006-01-02T15:04:05"))
	}
	dateTo, errPrsTo := time.ParseInLocation("2006-01-02T15:04:05", sb.SoapRq.Req.ToTime, time.Local)
	if errPrsTo != nil {
		dateTo = time.Date(2038, 01, 19, 3, 14, 7, 0, time.UTC)
		logger.Errorf("Date `To` parsing error: %v; setting default: %v", err, dateTo.Format("2006-01-02T15:04:05"))
	}

	foundAcc, err := service.GetAcctInfoG2b(sb.SoapRq.Req.Account)
	if err != nil {
		return nil, err
	}

	cards, err := service.GetCardsListG2b(foundAcc.Custcode, foundAcc.Currcode)
	if err != nil {
		return nil, err
	}
	var cardRows []CardRow
	for _, v := range cards {
		var cardPan string
		switch v.PAN {
		case "":
			cardPan = v.LkeyDisplay
		default:
			cardPan = v.PAN
		}
		cardRows = append(cardRows, CardRow{
			PAN:    cardPan,
			MBR:    "0",
			Status: utils.CardStatuses[v.StatCode],
			Type:   utils.CardTypes[v.ProductType],
		},
		)
	}

	var (
		size     int
		cardTrns *d8corp.CardInfoData
	)
	size, err = strconv.Atoi(sb.SoapRq.Req.Count)
	if err != nil {
		logger.Errorf("[SERVICE] getAcctStatement req error: wrong Count param! Setting default: 10;")
		size = 10
	}
	if len(cardRows) != 0 {
		cardTrns, err = service.GetCardTransactionHistory(cardRows[0].PAN, dateFrom.Format("20060102150405"), dateTo.Format("20060102150405"), size)
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

	for i, v := range cardTrns.CardTransactions {
		operDate, _ := time.ParseInLocation("20060102", v.BusDate, time.Local)
		tranTime, _ := time.ParseInLocation("20060102150405", v.When_created[:14], time.Local)

		resp.Statement.Rows = append(resp.Statement.Rows, Row{
			FrontId:             fmt.Sprintf("%d", v.TlId),
			Origin:              v.EcTxRefno, //"XAPI/00005GSLZE4o2Ddn2iCG7KrtMmnxg5Va",
			Type:                "1",         // 1–финансовая, 2–внутридоговорная, 3–авторизационная
			OperCode:            fmt.Sprintf("%d", v.Txncode),
			Description:         fmt.Sprintf("%s %s", v.CrdactplocName, v.EcTxRefno), //"*Dushanbe RRP XAPI/00005GSLZE4o2Ddn2iCG7KrtMmnxg5Va",
			Amount:              fmt.Sprintf("%.2f", v.TxnAmount),
			OperDate:            operDate.Format("2006-01-02"),
			TranTime:            tranTime.Format("2006-01-02"),
			OrigAmount:          fmt.Sprintf("%.2f", v.Amtbill),
			OrigCurrency:        utils.Currencies[v.Curbill],
			PAN:                 v.Lkey.Pan,
			MBR:                 "0",
			TermClass:           v.TermType,
			TermName:            v.TermCode,
			TermSIC:             fmt.Sprintf("%d", v.CrdacptBus),
			TermLocation:        v.CrdactplocName,
			ApprovalCode:        fmt.Sprintf("%.6d", v.Stan/1000000),
			SeqNo:               fmt.Sprintf("%v", i),
			TermCountry:         v.CrdactplocCountry,
			TermCity:            v.CrdactplocCity,
			OnlineIssuerFee:     "0",
			OrigTime:            tranTime.Format("2006-01-02T15:04:05"),
			Currency:            utils.Currencies[v.TxnCurrency],
			TermRetailerName:    v.CrdactplocName,
			CurrencyISOCode:     v.TxnCurrency,
			OrigCurrencyISOCode: v.Curbill,
		})
		resp.Statement.Rows[i].SeqNo = fmt.Sprintf("%d", i+1)
	}

	soapResp.Body = RespBody{
		GetAcctStatementRp: GetAcctStatementRp{
			Response: resp,
		},
	}
	return soapResp, nil
}
