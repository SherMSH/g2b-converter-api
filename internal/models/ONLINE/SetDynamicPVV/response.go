package setdynamicpvv

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
	SetDynamicPVVRp SetDynamicPVVRp `xml:"m1:SetDynamicPVV_PINOffsetRp"`
}

type SetDynamicPVVRp struct {
	Response Response `xml:"m1:Response"`
}

type Response struct {
	Echo         string `xml:"Echo,attr"`
	Product      string `xml:"Product,attr"`
	ResponseAttr string `xml:"Response,attr"`
	TranId       string `xml:"TranId,attr"`
	Ver          string `xml:"Ver,attr"`
}
