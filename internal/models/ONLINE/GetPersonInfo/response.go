package getpersoninfo

import "encoding/xml"

type Envelope struct {
	XMLName xml.Name `xml:"s:Envelope"`
	XmlnsS  string   `xml:"xmlns:s,attr"`
	XmlnsM1 string   `xml:"xmlns:m1,attr"`
	XmlnsM0 string   `xml:"xmlns:m0,attr"`
	Body    RespBody `xml:"s:Body"`
}

type RespBody struct {
	GetPersonInfoRp GetPersonInfoRp `xml:"m1:GetPersonInfoRp"`
}

type GetPersonInfoRp struct {
	Response Response `xml:"m1:Response"`
}

type Response struct {
	Echo         string `xml:"Echo,attr"`
	Product      string `xml:"Product,attr"`
	ResponseAttr string `xml:"Response,attr"`
	TranId       string `xml:"TranId,attr"`
	Ver          string `xml:"Ver,attr"`

	Accounts             []Account            `xml:"Accounts>Row"`
	AlternativeMessaging []AlternativeMessage `xml:"AlternativeMessaging>Row"`
	Cards                []Card               `xml:"Cards>Row"`
	Confidential         []Confidential       `xml:"Confidential>Row"`
	Info                 []Info               `xml:"Info>Row"`
}

// Аккаунт/счет
type Account struct {
	PersonId string `xml:"PersonId"`
	Account  string `xml:"Account"`
	Type     string `xml:"Type"`
	Status   string `xml:"Status"`
}

// Альтернативный канал связи (SMS и т.д.)
type AlternativeMessage struct {
	Channel        string `xml:"Channel"`
	Address        string `xml:"Address"`
	Scheme         string `xml:"Scheme"`
	Disabled       string `xml:"Disabled"`
	UseForDynAuth  string `xml:"UseForDynAuth"`
	IsDefault      string `xml:"IsDefault"`
	Broadcast      string `xml:"Broadcast"`
	DynAuthBlocked string `xml:"DynAuthBlocked"`
	ErrorCount     string `xml:"ErrorCount"`
	CreationDate   string `xml:"CreationDate"`
	LastUpdateTime string `xml:"LastUpdateTime"`
}

// Банковская карта
type Card struct {
	PersonId string `xml:"PersonId"`
	PAN      string `xml:"PAN"`
	MBR      string `xml:"MBR"`
	Status   string `xml:"Status"`
	ExpDate  string `xml:"ExpDate"`
	Type     string `xml:"Type"`
}

// Конфиденциальная информация
type Confidential struct {
	PersonId     string `xml:"PersonId"`
	What         string `xml:"What"`
	Value        string `xml:"Value"`
	IsAllowedCST string `xml:"IsAllowedCST"`
	IsAllowedADS string `xml:"IsAllowedADS"`
	IsAllowedTB  string `xml:"IsAllowedTB"`
}

// Персональная информация
type Info struct {
	PersonId        string `xml:"PersonId"`
	FIO             string `xml:"FIO"`
	Sex             string `xml:"Sex"`
	IdentType       string `xml:"IdentType"`
	Identity        string `xml:"Identity"`
	Birthday        string `xml:"Birthday"`
	Birthplace      string `xml:"Birthplace"`
	VIP             string `xml:"VIP"`
	InstName        string `xml:"InstName"`
	PersonExtId     string `xml:"PersonExtId"`
	ResidentCountry string `xml:"ResidentCountry"`
	AddressInLatin  string `xml:"AddressInLatin"`
	FirstNameNat    string `xml:"FirstNameNat"`
	LastNameNat     string `xml:"LastNameNat"`
	MiddleNameNat   string `xml:"MiddleNameNat"`
	TaxPayerNumber  string `xml:"TaxPayerNumber"`
	AddressNat      string `xml:"AddressNat"`
}
