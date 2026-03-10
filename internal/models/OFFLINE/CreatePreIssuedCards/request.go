package createpreissuedcards

import "encoding/xml"

// Root - корневой элемент XML
type Root struct {
	XMLName xml.Name `xml:"ROOT"`
	Record  []Record `xml:"RECORD"`
}

// Record - основная запись с данными клиента
type Record struct {
	// Персональные данные
	FIO        string `xml:"FIO"`
	SEX        string `xml:"SEX"`
	Title      string `xml:"TITLE"`
	LatFIO     string `xml:"LATFIO"`
	BirthFIO   string `xml:"BIRTHFIO"`
	BirthPlace string `xml:"BIRTHPLACE"`
	Family     string `xml:"FAMILY"`
	Education  string `xml:"EDUCATION"`
	Occupation string `xml:"OCCUPATION"`
	Birthday   string `xml:"BIRTHDAY"`

	// Паспортные данные
	PasNom    string `xml:"PASNOM"`
	PasDat    string `xml:"PASDAT"`
	PasExpDat string `xml:"PASEXPDAT"`
	PasPlace  string `xml:"PASPLACE"`
	PasDep    string `xml:"PASDEP"`

	// Резидентство
	Resident   string `xml:"RESIDENT"`
	CountryRes string `xml:"COUNTRYRES"`
	ExtID      string `xml:"EXTID"`
	PCode      string `xml:"PCODE"`
	BRPart     string `xml:"BRPART"`
	UserData   string `xml:"USERDATA"`
	STLang     string `xml:"STLANG"`
	StartBank  string `xml:"STARTBANK"`
	VIP        string `xml:"VIP"`

	// Секретная информация
	SecretInfo SecretInfo `xml:"SECRETINFO"`

	// Разрешения
	IsAllowedCST string `xml:"ISALLOWEDCST"`
	IsAllowedADS string `xml:"ISALLOWEDADS"`
	IsAllowedTBU string `xml:"ISALLOWEDTBU"`

	// ИНН и адрес проживания
	INN        string `xml:"INN"`
	Address    string `xml:"ADDRESS"`
	ZipLive    string `xml:"ZIPLIVE"`
	CntrLive   string `xml:"CNTRYLIVE"`
	RegionLive string `xml:"REGIONLIVE"`
	CityLive   string `xml:"CITYLIVE"`
	StreetLive string `xml:"STREETLIVE"`
	HouseLive  string `xml:"HOUSELIVE"`
	BuildLive  string `xml:"BUILDLIVE"`
	FrameLive  string `xml:"FRAMELIVE"`
	FlatLive   string `xml:"FLATLIVE"`

	// Адрес регистрации
	ResAddress string `xml:"RESADDRESS"`
	ZipReg     string `xml:"ZIPREG"`
	CntrReg    string `xml:"CNTRYREG"`
	RegionReg  string `xml:"REGIONREG"`
	CityReg    string `xml:"CITYREG"`
	StreetReg  string `xml:"STREETREG"`
	HouseReg   string `xml:"HOUSEREG"`
	BuildReg   string `xml:"BUILDREG"`
	FrameReg   string `xml:"FRAMEREG"`
	FlatReg    string `xml:"FLATREG"`

	// Контактный адрес
	CorAddress string `xml:"CORADDRESS"`
	ZipCont    string `xml:"ZIPCONT"`
	CntrCont   string `xml:"CNTRYCONT"`
	RegionCont string `xml:"REGIONCONT"`
	CityCont   string `xml:"CITYCONT"`
	StreetCont string `xml:"STREETCONT"`
	HouseCont  string `xml:"HOUSECONT"`
	BuildCont  string `xml:"BUILDCONT"`
	FrameCont  string `xml:"FRAMECONT"`
	FlatCont   string `xml:"FLATCONT"`

	// Контакты
	Email     string `xml:"EMAIL"`
	Fax       string `xml:"FAX"`
	Phone     string `xml:"PHONE"`
	CellPhone string `xml:"CELLPHONE"`
	Pager     string `xml:"PAGER"`

	// Работа
	Company  string `xml:"COMPANY"`
	Ceh      string `xml:"CEH"`
	TabNom   string `xml:"TABNOM"`
	StartJob string `xml:"STARTJOB"`
	Job      string `xml:"JOB"`
	JobPhone string `xml:"JOBPHONE"`
	Salary   string `xml:"SALARY"`

	// Банковские данные
	Affiliate  string `xml:"AFFILIATE"`
	Account    string `xml:"ACCOUNT"`
	ExtAccount string `xml:"EXTACCOUNT"`
	AcctType   string `xml:"ACCTTYPE"`
	AcctStat   string `xml:"ACCTSTAT"`
	AccountTP  string `xml:"ACCOUNTTP"`
	Acct2CDesc string `xml:"ACCT2CDESC"`
	AccFinProf string `xml:"ACCFINPROF"`
	ExternalID string `xml:"EXTERNALID"`
	CancelDate string `xml:"CANCELDATE"`
	MakePrior  string `xml:"MAKEPRIOR"`
	OrderDate  string `xml:"ORDERDATE"`
	NameOnCard string `xml:"NAMEONCARD"`
	CardPrefix string `xml:"CARDPREFIX"`
	CRDStat    string `xml:"CRDSTAT"`
	GroupCMD   string `xml:"GROUPCMD"`
	FinProf    string `xml:"FINPROF"`
	FinProfExt string `xml:"FINPROFEXT"`
	CurrencyNo string `xml:"CURRENCYNO"`
	RiskLevel  string `xml:"RISKLEVEL"`
	PinOffset  string `xml:"PINOFFSET"`
	CVV        string `xml:"CVV"`
	CVV2       string `xml:"CVV2"`
	ECStatus   string `xml:"ECSTATUS"`
	CRDProdID  string `xml:"CRDPRODID"`
	PAN        string `xml:"PAN"`
	MBR        string `xml:"MBR"`
}

// SecretInfo - элемент с секретной информацией
type SecretInfo struct {
	Items []SecretItem `xml:"item"`
}

// SecretItem - элемент item внутри SECRETINFO
type SecretItem struct {
	Ind   string `xml:"ind,attr"`
	What  string `xml:"what,attr"`
	Value string `xml:"value,attr"`
}
