package getacctstatement

import (
	d8corp "converterapi/internal/models/D8CORP"
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	d8procweb "converterapi/pkg/d8-proc-web"
	"converterapi/pkg/logger"
	"fmt"
	"strconv"
	"time"
)

const defaultStatementSize = 10

// Svc собирает выписку по счёту.
//
// Источник данных - xmiss/getAccountStatement (7.30): он отдаёт движения самого
// счёта, включая операции, прошедшие мимо карты. Строки, связанные с
// транзакциями, дополняются реквизитами из getTransactionDetails - в выписке
// счёта нет ни кода авторизации, ни данных терминала.
func Svc(sb *Body) (soapResp *Envelope, err error) {
	d8procweb.Signin()
	defer d8procweb.Signout()

	dateFrom, errPrsFrom := time.ParseInLocation("2006-01-02T15:04:05", sb.SoapRq.Req.FromTime, time.Local)
	if errPrsFrom != nil {
		dateFrom = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		logger.Errorf("Date `From` parsing error: %v; setting default: %v", errPrsFrom, dateFrom.Format("2006-01-02T15:04:05"))
	}
	dateTo, errPrsTo := time.ParseInLocation("2006-01-02T15:04:05", sb.SoapRq.Req.ToTime, time.Local)
	if errPrsTo != nil {
		dateTo = time.Date(2038, 1, 19, 3, 14, 7, 0, time.UTC)
		logger.Errorf("Date `To` parsing error: %v; setting default: %v", errPrsTo, dateTo.Format("2006-01-02T15:04:05"))
	}

	size, err := strconv.Atoi(sb.SoapRq.Req.Count)
	if err != nil || size <= 0 {
		logger.Errorf("[SERVICE] getAcctStatement req error: wrong Count param! Setting default: %d;", defaultStatementSize)
		size = defaultStatementSize
	}

	// Счёт ищем, чтобы получить валюту: ключ getAccountStatement - номер + валюта
	foundAcc, err := service.GetAcctInfoG2b(sb.SoapRq.Req.Account)
	if err != nil {
		return nil, err
	}

	statement, err := service.GetAccountStatementG2b(
		foundAcc.Accnum,
		foundAcc.Currcode,
		dateFrom.Format("20060102"),
		dateTo.Format("20060102"),
		size, 1,
	)
	if err != nil {
		return nil, err
	}

	// Детали нужны только строкам, порождённым транзакциями
	refs := make([]service.TxRef, 0, len(statement.AccountLog))
	for _, rec := range statement.AccountLog {
		if rec.TlId != 0 {
			refs = append(refs, service.TxRef{TlId: rec.TlId})
		}
	}
	details := service.GetTransactionDetailsBatch(refs)

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
	resp.Statement.Rows = make([]Row, 0, len(statement.AccountLog))

	for i, rec := range statement.AccountLog {
		row := Row{
			FrontId:         frontId(rec),
			Type:            "1",
			Description:     rec.Description,
			Origin:          rec.ExtTxnId,
			Amount:          fmt.Sprintf("%.2f", rec.NewState.AvlBal-rec.OldState.AvlBal),
			Remain:          fmt.Sprintf("%.2f", rec.NewState.AvlBal),
			OperDate:        operDate(rec.DateLocal),
			TranTime:        tstamp(rec.TstampInsert),
			OrigTime:        tstamp(rec.TstampInsert),
			Currency:        currency(rec.NewState.Currency),
			CurrencyISOCode: rec.NewState.Currency,
			MBR:             "0",
			OnlineIssuerFee: "0",
			SeqNo:           strconv.Itoa(i + 1),
		}

		if trn, ok := details[rec.TlId]; ok {
			enrich(&row, trn)
		}
		resp.Statement.Rows = append(resp.Statement.Rows, row)
	}

	soapResp.Body = RespBody{
		GetAcctStatementRp: GetAcctStatementRp{
			Response: resp,
		},
	}
	return soapResp, nil
}

// enrich дополняет строку выписки реквизитами транзакции: их нет в движении по счёту
func enrich(row *Row, trn d8corp.TransactionDetails) {
	row.OperCode = utils.TranCode(trn.TxnCode)
	row.ApprovalCode = trn.Aprvlcode
	row.PAN = trn.Lkey.Pan
	row.TermClass = trn.Termtype
	row.TermName = trn.TermCode
	row.TermSIC = strconv.Itoa(trn.CrdacptBus)
	row.TermLocation = trn.CrdacptlocName
	row.TermRetailerName = trn.CrdacptlocName
	row.TermCity = trn.CrdacptlocCity
	row.TermCountry = trn.CrdacptlocCountry
	row.OrigAmount = fmt.Sprintf("%.2f", trn.Amtbill)
	row.OrigCurrency = currency(trn.Curbill)
	row.OrigCurrencyISOCode = trn.Curbill
	if trn.EcTxRefno != "" {
		row.Origin = trn.EcTxRefno
	}
	if row.Description == "" {
		row.Description = trn.CrdacptlocName
	}
}

// frontId - идентификатор строки: транзакция, если движение её породило,
// иначе идентификатор записи журнала счёта
func frontId(rec d8corp.AccountLog) string {
	if rec.TlId != 0 {
		return strconv.Itoa(rec.TlId)
	}
	return strconv.FormatInt(rec.Id, 10)
}

// currency переводит числовой код в буквенный, не подставляя заглушку "unknown"
func currency(code string) string {
	if code == "" {
		return ""
	}
	return utils.Currencies[code]
}

func operDate(dateLocal string) string {
	t, err := time.ParseInLocation("20060102", dateLocal, time.Local)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func tstamp(tstampInsert string) string {
	if len(tstampInsert) < 14 {
		return ""
	}
	t, err := time.ParseInLocation("20060102150405", tstampInsert[:14], time.Local)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02T15:04:05")
}
