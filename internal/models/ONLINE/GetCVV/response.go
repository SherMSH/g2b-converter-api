package getcvv

import "encoding/xml"

type Envelope struct {
	XMLName xml.Name `xml:"s:Envelope"`
	XmlnsS  string   `xml:"xmlns:s,attr"`
	XmlnsM1 string   `xml:"xmlns:m1,attr"`
	XmlnsM0 string   `xml:"xmlns:m0,attr"`
	Body    RespBody `xml:"s:Body"`
}

type RespBody struct {
	GetCVVRp GetCVVRp `xml:"m1:GetCVVRp"`
}

type GetCVVRp struct {
	Response Response `xml:"m1:Response"`
}

type Response struct {
	Echo         string `xml:"Echo,attr"`
	Product      string `xml:"Product,attr"`
	ResponseAttr string `xml:"Response,attr"`
	TranId       string `xml:"TranId,attr"`
	Ver          string `xml:"Ver,attr"`

	CVV                  string `xml:"m0:CVV"`
	CardVerificationType string `xml:"m0:CardVerificationType"`
	StrCVV               string `xml:"m0:StrCVV"`
}
