package relinkpreissuedcards

import "encoding/xml"

type Envelope struct {
	XMLName xml.Name `xml:"s:Envelope"`
	XmlnsS  string   `xml:"xmlns:s,attr"`
	XmlnsM1 string   `xml:"xmlns:m1,attr"`
	XmlnsM0 string   `xml:"xmlns:m0,attr"`
	Body    RespBody `xml:"s:Body"`
}

type RespBody struct {
	UpdateCard2AcctLinkRp UpdateCard2AcctLinkRp `xml:"m1:UpdateCard2AcctLinkRp"`
	DeleteCard2AcctLinkRp DeleteCard2AcctLinkRp `xml:"m1:DeleteCard2AcctLinkRp"`
	SetCardPersonRp       SetCardPersonRp       `xml:"m1:SetCardPersonRp"`
	SetCardStatusRp       SetCardStatusRp       `xml:"m1:SetCardStatusRp"`
}

type UpdateCard2AcctLinkRp struct {
	Response Response `xml:"m1:Response"`
}

type DeleteCard2AcctLinkRp struct {
	Response Response `xml:"m1:Response"`
}

type SetCardPersonRp struct {
	Response Response `xml:"m1:Response"`
}

type SetCardStatusRp struct {
	Response Response `xml:"m1:Response"`
}

type Response struct {
	Echo         string `xml:"Echo,attr"`
	Product      string `xml:"Product,attr"`
	ResponseAttr string `xml:"Response,attr"`
	TranId       string `xml:"TranId,attr"`
	Ver          string `xml:"Ver,attr"`
}
