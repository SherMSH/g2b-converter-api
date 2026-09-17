package verifypin

import "reflect"

type Body struct {
	SoapRq SoapRq `xml:"VerifyPINRq" json:"VerifyPINRq"`
}

func (sb Body) GetBodyType() reflect.Type {
	return reflect.TypeOf(sb)
}

func (sb *Body) Call() (*Envelope, error) {
	return Svc(sb)
}

// SoapRq соответствует элементу fimi:VerifyPINRq
type SoapRq struct {
	Req Request `xml:"Request" json:"Request"`
}

// Request соответствует элементу fimi:Request
type Request struct {
	Ver      string `xml:"Ver,attr" json:"ver"`
	Product  string `xml:"Product,attr" json:"product"`
	Echo     string `xml:"Echo,attr" json:"echo"`
	Session  string `xml:"Session,attr" json:"session"`
	Clerk    string `xml:"Clerk,attr" json:"clerk"`
	Password string `xml:"Password,attr" json:"password"`

	PAN string `xml:"PAN" json:"pan"`
	MBR string `xml:"MBR" json:"mbr"`
	PIN string `xml:"PIN" json:"pin"`

	// Срок действия в формате YYMM. Если не прислан - берём из карты
	ExpirationDate string `xml:"ExpirationDate" json:"expiration_date"`

	PersonId string `xml:"PersonId" json:"person_id"`
	CardUID  string `xml:"CardUID" json:"card_uid"`
}
