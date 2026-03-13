package addpersoncmsabonent

import "reflect"

type Body struct {
	SoapRq SoapRq `xml:"AddPersonCMSAbonentRq" json:"AddPersonCMSAbonentRq"`
}

func (sb Body) GetBodyType() reflect.Type {
	return reflect.TypeOf(sb)
}

// SoapRq соответствует элементу AddCMSAbonentRq
type SoapRq struct {
	Req Request `xml:"Request" json:"request"`
}

// Request соответствует элементу Request
type Request struct {
	Ver      string `xml:"Ver,attr" json:"ver"`
	Product  string `xml:"Product,attr" json:"product"`
	Echo     string `xml:"Echo,attr" json:"echo"`
	Clerk    string `xml:"Clerk,attr" json:"clerk"`
	Password string `xml:"Password,attr" json:"password"`

	InstName     string       `xml:"InstName" json:"inst_name"`
	PersonId     string       `xml:"PersonId" json:"person_id"`
	AltMessaging AltMessaging `xml:"AlternativeMessaging" json:"alternative_messaging"`
	NeedNotify   string       `xml:"NeedNotify" json:"need_notify"`
}

type AltMessaging struct {
	Row []Row `xml:"Row" json:"rows"`
}

type Row struct {
	Channel       string `xml:"Channel" json:"channel"`
	Address       string `xml:"Address" json:"address"`
	Scheme        string `xml:"Scheme" json:"scheme"`
	Disabled      string `xml:"Disabled" json:"disabled"`
	UseForDynAuth string `xml:"UseForDynAuth" json:"use_for_dyn_auth"`
	IsDefault     string `xml:"IsDefault" json:"is_default"`
	Priority      string `xml:"Priority" json:"priority"`
	Broadcast     string `xml:"Broadcast" json:"broadcast"`

	// Name  string `xml:"Name" json:"name"`
	// Title string `xml:"Title" json:"title"`
}
