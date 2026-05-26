package getcvv

import "reflect"

type Body struct {
	SoapRq SoapRq `xml:"GetCVVRq" json:"GetCVVRq"`
}

func (sb Body) GetBodyType() reflect.Type {
	return reflect.TypeOf(sb)
}

func (sb *Body) Call() (*Envelope, error) {
	return Svc(sb)
}

// SoapRq соответствует элементу GetAcctInfoRq
type SoapRq struct {
	Req Request `xml:"Request" json:"request"`
}

// Request соответствует элементу Request
type Request struct {
	Ver      string `xml:"Ver,attr" json:"ver"`
	Product  string `xml:"Product,attr" json:"product"`
	Echo     string `xml:"Echo,attr" json:"echo"`
	Clerk    string `xml:"Clerk,attr" json:"clerk"`
	Password string `xml:"Password,attr" json:"password"`

	PAN     string `xml:"PAN" json:"pan"`
	MBR     string `xml:"MBR" json:"mbr"`
	CardUID string `xml:"CardUID" json:"card_uid"`
	IsCVV2  string `xml:"IsCVV2" json:"is_cvv2"`
	ExpDate string `xml:"ExpDate" json:"exp_date"`
}
