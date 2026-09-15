package posrequestrq

import "encoding/xml"

// Envelope - корневой элемент SOAP конверта
type Envelope struct {
	XMLName xml.Name `xml:"s:Envelope"`
	XmlnsS  string   `xml:"xmlns:s,attr"`
	XmlnsM1 string   `xml:"xmlns:m1,attr"`
	XmlnsM0 string   `xml:"xmlns:m0,attr"`
	Body    RespBody `xml:"s:Body"`
}

// Body - тело SOAP сообщения
type RespBody struct {
	POSRequestRp POSRequestRp `xml:"m1:POSRequestRp"`
}

// POSRequestRp - корневой элемент запроса
type POSRequestRp struct {
	Response Response `xml:"m1:Response"`
}

// Response - основной элемент ответа
type Response struct {
	Echo         string `xml:"Echo,attr"`
	Product      string `xml:"Product,attr"`
	ResponseAttr string `xml:"Response,attr"`
	TranId       string `xml:"TranId,attr"`
	Ver          string `xml:"Ver,attr"`

	AccountCurrency      string      `xml:"m0:AccountCurrency,omitempty"`
	ApprovalCode         string      `xml:"m0:ApprovalCode,omitempty"`
	AuthRespCode         string      `xml:"m0:AuthRespCode"`
	AuthRespCodeCategory string      `xml:"m0:AuthRespCodeCategory"`
	AvailBalance         string      `xml:"m0:AvailBalance,omitempty"`
	BalanceCurrency      string      `xml:"m0:BalanceCurrency"`
	BonusDebt            string      `xml:"m0:BonusDebt"`
	CVxOK                string      `xml:"m0:CVxOK"`
	Currency             string      `xml:"m0:Currency"`
	DeclineReason        string      `xml:"m0:DeclineReason,omitempty"`
	Fee                  string      `xml:"m0:Fee"`
	FromAcct             string      `xml:"m0:FromAcct,omitempty"`
	IssuerFee            string      `xml:"m0:IssuerFee,omitempty"`
	LedgerBalance        string      `xml:"m0:LedgerBalance,omitempty"`
	MaskBalances         string      `xml:"m0:MaskBalances"`
	RelatedTran          RelatedTran `xml:"m0:RelatedTran"`
	ThisTranId           string      `xml:"m0:ThisTranId"`
	ToAcct               string      `xml:"m0:ToAcct,omitempty"`
}

type RelatedTran struct {
	Rows []Rows `xml:"m0:Row"`
}
type Rows struct {
	RelatedTranId       string `xml:"m0:RelatedTranId"`
	RelatedTranCode     string `xml:"m0:RelatedTranCode"`
	RelatedAuthRespCode string `xml:"m0:RelatedAuthRespCode"`
}
