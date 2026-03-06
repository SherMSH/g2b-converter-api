package updateperson

import "reflect"

type Body struct {
	SoapRq SoapRq `xml:"UpdatePersonRq" json:"UpdatePersonRq"`
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

	Id                 string `xml:"Id" json:"id"`
	InstName           string `xml:"InstName" json:"inst_name"`
	FIO                string `xml:"FIO" json:"fio"`
	IdentType          string `xml:"IdentType" json:"ident_type"`
	UpdateFieldsMethod string `xml:"UpdateFieldsMethod" json:"update_fields_method"`
	FirstName          string `xml:"FirstName" json:"first_name"`
	LastName           string `xml:"LastName" json:"last_name"`
	MiddleName         string `xml:"MiddleName" json:"middle_name"`

	BranchId                string `xml:"BranchId" json:"branch_id"`
	VIP                     string `xml:"VIP" json:"vip"`
	Sex                     string `xml:"Sex" json:"sex"`
	Identity                string `xml:"Identity" json:"identity"`
	SecretWord              string `xml:"SecretWord" json:"secret_word"`
	Birthday                string `xml:"Birthday" json:"birthday"`
	Birthplace              string `xml:"Birthplace" json:"birthplace"`
	ResidentCountry         string `xml:"ResidentCountry" json:"resident_country"`
	ResidentCityInLatin     string `xml:"ResidentCityInLatin" json:"resident_country_in_latin"`
	AddressInLatin          string `xml:"AddressInLatin" json:"address_in_latin"`
	PostalCode              string `xml:"PostalCode" json:"PostalCode"`
	ParentPersonId          string `xml:"ParentPersonId" json:"parent_person_id"`
	ParentInstName          string `xml:"ParentInstName" json:"parent_inst_name"`
	ChangeReason            string `xml:"ChangeReason" json:"change_reason"`
	ResidentState           string `xml:"ResidentState" json:"resident_state"`
	ResidentStateExternalId string `xml:"ResidentStateExternalId" json:"resident_state_external_id"`
	CategoryId              string `xml:"CategoryId" json:"category_id"`
}
