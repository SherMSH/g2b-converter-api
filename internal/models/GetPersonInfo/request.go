package getpersoninfo

import "reflect"

type Body struct {
	SoapRq SoapRq `xml:"GetPersonInfoRq" json:"GetPersonInfoRq"`
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

	FIO            string `xml:"FIO" json:"fio"`
	Id             string `xml:"Id" json:"id"`
	InstName       string `xml:"InstName" json:"inst_name"`
	Identity       string `xml:"Identity" json:"identity"`
	Birthday       string `xml:"Birthday" json:"birthday"`
	IdentType      string `xml:"IdentType" json:"ident_type"`
	AbonentAddress string `xml:"AbonentAddress" json:"abonent_address"`
	PersonExtId    string `xml:"PersonExtId" json:"person_ext_id"`
}
