package getcardinfo

import (
	"reflect"
)

type Body struct {
	SoapRq SoapRq `xml:"GetCardInfoRq" json:"GetCardInfoRq"`
}

func (sb Body) GetBodyType() reflect.Type {
	return reflect.TypeOf(sb)
}

func (sb *Body) Call() (*Envelope, error) {
	rsp, err := Svc(sb)
	return rsp, err
}

// SoapRq соответствует элементу GetCardInfoRq
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

	PAN          string `xml:"PAN" json:"pan"`
	MBR          string `xml:"MBR" json:"mbr"`
	Type         string `xml:"Type" json:"type"`
	PersonId     string `xml:"PersonId" json:"person_id"`
	CardUID      string `xml:"CardUID" json:"card_uid"`
	RequiredData string `xml:"RequiredData" json:"required_data"`

	Status Status `xml:"Status" json:"status"`

	ExpirationDate string `xml:"ExpirationDate" json:"expiration_date"`
	NameOnCard     string `xml:"NameOnCard" json:"name_on_card"`
}

type Status struct {
	Row []Row `xml:"Row" json:"rows"`
}

type Row struct {
	Value string `xml:"Value" json:"value"`
}
