package initsession

import "reflect"

type Body struct {
	SoapRq SoapRq `xml:"InitSessionRq" json:"InitSessionRq"`
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

	NeedDicts string `xml:"NeedDicts" json:"need_dicts"`

	AllVendors    string `xml:"AllVendors" json:"all_vendors"`
	AvoidSession  string `xml:"AvoidSession" json:"avoid_session"`
	CurrencyCodes Rows   `xml:"CurrencyCodes" json:"currency_codes"`
}

type Rows struct {
	Row []Row `xml:"Row" json:"rows"`
}

type Row struct {
	Value string `xml:"Value" json:"value"`
}
