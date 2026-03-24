package createorganizations

import (
	"converterapi/internal/utils"
	"encoding/xml"
)

// Root - корневой элемент XML
type Root struct {
	XMLName xml.Name `xml:"ROOT"`
	Record  []Record `xml:"RECORD"`
}

func (r Root) GetReqType() utils.OfflineReqType {
	return utils.CreateOrganizations
}

func (r Root) Call() error {
	return nil
}

// Record - запись с данными о филиале
type Record struct {
	// Основная информация
	Name      string `xml:"NAME"`
	ExtID     string `xml:"EXTID"`
	ShortName string `xml:"SHORTNAME"`
	TransName string `xml:"TRANSNAME"`

	// Контактные данные
	PhoneLegal string `xml:"PHONELEGAL"`
	FaxLegal   string `xml:"FAXLEGAL"`
	EmailLegal string `xml:"EMAILLEGAL"`
	PhoneCont  string `xml:"PHONECONT"`
	FaxCont    string `xml:"FAXCONT"`
	EmailCont  string `xml:"EMAILCONT"`

	// Регистрационные данные
	RegNo    string `xml:"REGNO"`
	RegDate  string `xml:"REGDATE"`
	RegLoc   string `xml:"REGLOC"`
	TPN      string `xml:"TPN"`
	OwnType  string `xml:"OWNTYPE"`
	EntType  string `xml:"ENTTYPE"`
	ComCode  string `xml:"COMCODE"`
	StatDir  string `xml:"STATDIR"`
	StatLang string `xml:"STATLANG"`
	AddInfo  string `xml:"ADDINFO"`
	BRPart   string `xml:"BRPART"`

	// Адреса
	Legal   Address `xml:"LEGAL"`
	Contact Address `xml:"CONTACT"`
	Post    Address `xml:"POST"`

	// Банковская информация
	BankInfoList BankInfoList `xml:"BANKINFOLIST"`

	// Контактные лица
	ContactPersonList ContactPersonList `xml:"CONTACTPERSONLIST"`

	// Подразделения
	SubdivisionList SubdivisionList `xml:"SUBDIVISIONLIST"`
}

// Address - структура для адреса
type Address struct {
	Country string `xml:"COUNTRY"`
	Region  string `xml:"REGION"`
	Area    string `xml:"AREA"`
	City    string `xml:"CITY"`
	Street  string `xml:"STREET"`
	Zip     string `xml:"ZIP"`
	House   string `xml:"HOUSE"`
	Build   string `xml:"BUILD"`
	Frame   string `xml:"FRAME"`
	Flat    string `xml:"FLAT"`
	Line    string `xml:"LINE"` // Дополнительная строка адреса
}

// BankInfoList - список банковской информации
type BankInfoList struct {
	BankInfos []BankInfo `xml:"BANKINFO"`
}

// BankInfo - банковская информация
type BankInfo struct {
	BIC        string `xml:"BIC"`
	CorAccount string `xml:"CORACCOUNT"`
	Account    string `xml:"ACCOUNT"`
}

// ContactPersonList - список контактных лиц
type ContactPersonList struct {
	ContactPersons []ContactPerson `xml:"CONTACTPERSON"`
}

// ContactPerson - контактное лицо
type ContactPerson struct {
	Index      string `xml:"INDEX"`
	PasNo      string `xml:"PASNO"`
	PersonName string `xml:"PERSONNAME"`
	Title      string `xml:"TITLE"`
	Info       string `xml:"INFO"`
}

// SubdivisionList - список подразделений
type SubdivisionList struct {
	Subdivisions []Subdivision `xml:"SUBDIVISION"`
}

// Subdivision - подразделение
type Subdivision struct {
	ID   string `xml:"ID"`
	Prop string `xml:"PROP"`
}
