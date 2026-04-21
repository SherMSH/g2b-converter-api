package resetbadpintries

import "reflect"

type Body struct {
	SoapRq SoapRq `xml:"ResetBadPINTriesRq" json:"ResetBadPINTriesRq"`
}

func (sb Body) GetBodyType() reflect.Type {
	return reflect.TypeOf(sb)
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

	IsInet       string `xml:"IsInet" json:"is_net"`
	MBR          string `xml:"MBR" json:"mbr"`
	PAN          string `xml:"PAN" json:"pan"`
	ChangeReason string `xml:"ChangeReason" json:"change_reason"`

	PersonId string `xml:"PersonId" json:"person_id"`
	CardUID  string `xml:"CardUID" json:"card_uid"`
}
