package getacctstatement

import "encoding/xml"

// Envelope - корневой элемент SOAP конверта
type Envelope struct {
	XMLName xml.Name `xml:"s:Envelope"`
	XmlnsS  string   `xml:"xmlns:s,attr"`
	XmlnsM1 string   `xml:"xmlns:m1,attr"`
	XmlnsM0 string   `xml:"xmlns:m0,attr"`
	Body    RespBody `xml:"s:Body"`
}

type RespBody struct {
	GetAcctStatementRq GetAcctStatementRq `xml:"m1:GetAcctStatementRq"`
}

type GetAcctStatementRq struct {
	Response Response `xml:"m1:Response"`
}

type Response struct {
	Echo         string `xml:"Echo,attr"`
	Product      string `xml:"Product,attr"`
	ResponseAttr string `xml:"Response,attr"`
	TranId       string `xml:"TranId,attr"`
	Ver          string `xml:"Ver,attr"`

	Statement Statement `xml:"m0:Statement"`
}

type Statement struct {
	Rows []Row `xml:"m0:Row"`
}

type Row struct {
	FrontId             string `xml:"m0:FrontId"`
	Origin              string `xml:"m0:Origin"`
	Type                string `xml:"m0:Type"`
	OperCode            string `xml:"m0:OperCode"`
	Description         string `xml:"m0:Description"`
	Amount              string `xml:"m0:Amount"`
	OperDate            string `xml:"m0:OperDate"`
	TranTime            string `xml:"m0:TranTime"`
	OrigAmount          string `xml:"m0:OrigAmount"`
	OrigCurrency        string `xml:"m0:OrigCurrency"`
	PAN                 string `xml:"m0:PAN"`
	MBR                 string `xml:"m0:MBR"`
	TermClass           string `xml:"m0:TermClass"`
	TermName            string `xml:"m0:TermName"`
	TermSIC             string `xml:"m0:TermSIC"`
	TermLocation        string `xml:"m0:TermLocation"`
	ApprovalCode        string `xml:"m0:ApprovalCode"`
	SeqNo               string `xml:"m0:SeqNo"`
	TermCountry         string `xml:"m0:TermCountry"`
	TermCity            string `xml:"m0:TermCity"`
	OnlineIssuerFee     string `xml:"m0:OnlineIssuerFee"`
	OrigTime            string `xml:"m0:OrigTime"`
	Currency            string `xml:"m0:Currency"`
	TermRetailerName    string `xml:"m0:TermRetailerName"`
	CurrencyISOCode     string `xml:"m0:CurrencyISOCode"`
	OrigCurrencyISOCode string `xml:"m0:OrigCurrencyISOCode"`
}
