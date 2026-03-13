package setcardstatus

import "reflect"

type Body struct {
	SoapRq SoapRq `xml:"SetCardStatusRq" json:"SetCardStatusRq"`
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

	MBR          string `xml:"MBR" json:"mbr"`
	PAN          string `xml:"PAN" json:"pan"`
	Status       string `xml:"Status" json:"status"`
	ChangeReason string `xml:"ChangeReason" json:"change_reason"`
	NeedNotify   string `xml:"NeedNotify" json:"need_notify"`

	PersonId      string `xml:"PersonId" json:"person_id"`
	CardUID       string `xml:"CardUID" json:"card_uid"`
	ParentMBR     string `xml:"ParentMBR" json:"parent_mbr"`
	ParentPAN     string `xml:"ParentPAN" json:"parent_pan"`
	ParentCardUID string `xml:"ParentCardUID" json:"parent_card_uid"`
	Reissued      string `xml:"Reissued" json:"reissued"`
}
