package acctcredit

import "encoding/xml"

type Envelope struct {
	XMLName xml.Name `xml:"s:Envelope"`
	XmlnsS  string   `xml:"xmlns:s,attr"`
	XmlnsM1 string   `xml:"xmlns:m1,attr"`
	XmlnsM0 string   `xml:"xmlns:m0,attr"`
	Body    RespBody `xml:"s:Body"`
}

type RespBody struct {
	AcctCreditRp AcctCreditRp `xml:"m1:AcctCreditRp"`
}

type AcctCreditRp struct {
	Response Response `xml:"m1:Response"`
}

type Response struct {
	Echo         string `xml:"Echo,attr"`
	Product      string `xml:"Product,attr"`
	ResponseAttr string `xml:"Response,attr"`
	TranId       string `xml:"TranId,attr"`
	Ver          string `xml:"Ver,attr"`

	ApprovalCode  string  `xml:"m0:ApprovalCode"`
	AvailBalance  float64 `xml:"m0:AvailBalance"`
	LedgerBalance float64 `xml:"m0:LedgerBalance"`
}
