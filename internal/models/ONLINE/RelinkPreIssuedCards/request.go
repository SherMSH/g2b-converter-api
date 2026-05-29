package relinkpreissuedcards

import (
	"encoding/xml"
	"reflect"
)

// SoapEnvelope соответствует корневому элементу <Envelope>
type SoapEnvelope struct {
	XMLName xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	Body    Body     `xml:"Body" json:"body"`
}

type Body struct {
	UpdateCard2AccLink SoapRq `xml:"UpdateCard2AcctLinkRq" json:"UpdateCard2AcctLinkRq"`
	DeleteCard2AccLink SoapRq `xml:"DeleteCard2AcctLinkRq" json:"DeleteCard2AcctLinkRq"`
	SetCardPersonRq    SoapRq `xml:"SetCardPersonRq" json:"SetCardPersonRq"`
	SetCardStatusRq    SoapRq `xml:"SetCardStatusRq" json:"SetCardStatusRq"`
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
	Echo     string `xml:"Echo,attr" json:"echo,omitempty"`

	Account        string `xml:"Account" json:"account,omitempty"`
	AcctStatus     string `xml:"AcctStatus" json:"acct_status,omitempty"`
	MBR            string `xml:"MBR" json:"mbr"`
	PAN            string `xml:"PAN" json:"pan"`
	ChangeReason   string `xml:"ChangeReason" json:"change_reason"`
	PersonId       string `xml:"PersonId" json:"person_id,omitempty"`
	Description    string `xml:"Description" json:"description,omitempty"`
	IgnoreChecking string `xml:"IgnoreChecking" json:"ignore_checking,omitempty"`
	NewPersonId    string `xml:"NewPersonId" json:"new_person_id,omitempty"`
	Status         string `xml:"Status" json:"status"`
	NeedNotify     string `xml:"NeedNotify" json:"need_notify,omitempty"`

	CardUID       string `xml:"CardUID" json:"card_uid"`
	AccountUID    string `xml:"AccountUID" json:"account_uid,omitempty"`
	ParentPAN     string `xml:"ParentPAN" json:"parent_pan,omitempty"`
	ParentMBR     string `xml:"ParentMBR" json:"parent_mbr,omitempty"`
	ParentCardUID string `xml:"ParentCardUID" json:"parent_card_uid,omitempty"`
	Reissued      string `xml:"Reissued" json:"reissued,omitempty"`
	ExpirityDate  string `xml:"ExpirityDate" json:"expirity_date"`
}
