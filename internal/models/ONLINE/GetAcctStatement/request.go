package getacctstatement

import (
	"reflect"
)

type Body struct {
	SoapRq SoapRq `xml:"GetAcctStatementRq" json:"GetAcctStatementRq"`
}

func (sb Body) GetBodyType() reflect.Type {
	return reflect.TypeOf(sb)
}

func (sb *Body) Call() (*Envelope, error) {
	return Svc(sb)
}

// SoapRq соответствует элементу fimi:GetAcctStatementRq
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

	Account    string `xml:"Account" json:"account"`
	InstName   string `xml:"InstName" json:"inst_name"`
	Count      string `xml:"Count" json:"count"`
	PersonId   string `xml:"PersonId" json:"person_id"`
	AccountUID string `xml:"AccountUID" json:"account_uid"`

	FromTime string `xml:"FromTime" json:"from_time"`
	ToTime   string `xml:"ToTime" json:"to_time"`
	Language string `xml:"Language" json:"language"`
}
