package getcardstatement

import "reflect"

type Body struct {
	SoapRq SoapRq `xml:"GetCardStatementRq" json:"GetCardStatementRq"`
}

func (sb Body) GetBodyType() reflect.Type {
	return reflect.TypeOf(sb)
}

// SoapRq соответствует элементу GetCardStatementRq
type SoapRq struct {
	Req Request `xml:"Request" json:"request"`
}

// Request соответствует элементу Request
type Request struct {
	Ver      string `xml:"Ver,attr" json:"ver"`
	Product  string `xml:"Product,attr" json:"product"`
	Clerk    string `xml:"Clerk,attr" json:"clerk"`
	Password string `xml:"Password,attr" json:"password"`

	Count string `xml:"Count" json:"count"`
	MBR   string `xml:"MBR" json:"mbr"`
	PAN   string `xml:"PAN" json:"pan"`

	FromTime string `xml:"FromTime" json:"from_time"`
	ToTime   string `xml:"ToTime" json:"to_time"`

	PersonId string `xml:"PersonId" json:"person_id"`
	CardUID  string `xml:"CardUID" json:"card_uid"`
	Language string `xml:"Language" json:"language"`
}
