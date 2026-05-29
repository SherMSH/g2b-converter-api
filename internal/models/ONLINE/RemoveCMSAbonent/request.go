package removecmsabonent

import "reflect"

type Body struct {
	SoapRq SoapRq `xml:"RemoveCMSAbonentRq" json:"RemoveCMSAbonentRq"`
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
	Clerk    string `xml:"Clerk,attr" json:"clerk"`
	Password string `xml:"Password,attr" json:"password"`
	Echo     string `xml:"Echo,attr" json:"echo"`

	MBR                  string `xml:"MBR" json:"mbr"`
	PAN                  string `xml:"PAN" json:"pan"`
	NeedNotify           string `xml:"NeedNotify" json:"need_notify"`
	AlternativeMessaging Rows   `xml:"AlternativeMessaging" json:"alternative_messaging"`

	CardUID      string `xml:"CardUID" json:"card_uid"`
	ChangeReason string `xml:"ChangeReason" json:"change_reason"`
}

type Rows struct {
	Row []Row `xml:"Row" json:"rows"`
}

type Row struct {
	Channel string `xml:"Channel" json:"channel"`
	Address string `xml:"Address" json:"address"`
}
