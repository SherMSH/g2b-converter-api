package getcardinfo

import "encoding/xml"

type Response struct {
	XMLName  xml.Name `xml:"Response"`
	Product  string   `xml:"Product,attr"`
	Response string   `xml:"Response,attr"`
	Ver      string   `xml:"Ver,attr"`

	Accounts                Accounts           `xml:"Accounts"`
	Acct2CardAttachType     string             `xml:"Acct2CardAttachType"`
	CNSDisabled             string             `xml:"CNSDisabled"`
	CardAllowedEMVScript    string             `xml:"CardAllowedEMVScript"`
	CardProfiles            CardProfiles       `xml:"CardProfiles"`
	ContactlessStatus       string             `xml:"ContactlessStatus"`
	ECNeedCAPAuth           string             `xml:"ECNeedCAPAuth"`
	ECNeedDynPwdAuth        string             `xml:"ECNeedDynPwdAuth"`
	ECNeedStaticAuth        string             `xml:"ECNeedStaticAuth"`
	ECNeedTokenAuth         string             `xml:"ECNeedTokenAuth"`
	ECStatus                string             `xml:"ECStatus"`
	ECUseCardSettingsAuth   string             `xml:"ECUseCardSettingsAuth"`
	ECUseDecoupledAuth      string             `xml:"ECUseDecoupledAuth"`
	EMVOptionsCheckDisabled string             `xml:"EMVOptionsCheckDisabled"`
	ExpDate                 string             `xml:"ExpDate"`
	FoundMBR                string             `xml:"FoundMBR"`
	FoundPAN                string             `xml:"FoundPAN"`
	IB_Registered           string             `xml:"IB_Registered"`
	InstName                string             `xml:"InstName"`
	IssueTechnology         string             `xml:"IssueTechnology"`
	LastATMUsed             string             `xml:"LastATMUsed"`
	LastChangeStatusTime    string             `xml:"LastChangeStatusTime"`
	LastPOSUsed             string             `xml:"LastPOSUsed"`
	LastPVVChangeTime       string             `xml:"LastPVVChangeTime"`
	LastRefreshTime         string             `xml:"LastRefreshTime"`
	LastTranId              string             `xml:"LastTranId"`
	LastTranTime            string             `xml:"LastTranTime"`
	MaskBalances            string             `xml:"MaskBalances"`
	MaskPVV                 string             `xml:"MaskPVV"`
	NameOnCard              string             `xml:"NameOnCard"`
	PINVerifyType           string             `xml:"PINVerifyType"`
	PVV                     string             `xml:"PVV"`
	PasswordFlag            string             `xml:"PasswordFlag"`
	PersonConfidential      PersonConfidential `xml:"PersonConfidential"`
	PersonExtId             string             `xml:"PersonExtId"`
	PersonFIO               string             `xml:"PersonFIO"`
	PersonId                string             `xml:"PersonId"`
	PersonVIP               string             `xml:"PersonVIP"`
	RequiredPasswordVersion string             `xml:"RequiredPasswordVersion"`
	RiskControlDisabled     string             `xml:"RiskControlDisabled"`
	RiskLevel               string             `xml:"RiskLevel"`
	Status                  string             `xml:"Status"`
	TmpECStatus             string             `xml:"TmpECStatus"`
	Type                    string             `xml:"Type"`
	UseUdCVV2               string             `xml:"UseUdCVV2"`
}

type Accounts struct {
	XMLName xml.Name   `xml:"Accounts"`
	Row     AccountRow `xml:"Row"`
}

type AccountRow struct {
	XMLName       xml.Name `xml:"Row"`
	AcctNo        string   `xml:"AcctNo"`
	Status        string   `xml:"Status"`
	LedgerBalance string   `xml:"LedgerBalance"`
	AvailBalance  string   `xml:"AvailBalance"`
	Currency      string   `xml:"Currency"`
	Type          string   `xml:"Type"`
	AccountStatus string   `xml:"AccountStatus"`
}

type CardProfiles struct {
	XMLName xml.Name       `xml:"CardProfiles"`
	Row     CardProfileRow `xml:"Row"`
}

type CardProfileRow struct {
	XMLName xml.Name `xml:"Row"`
	Id      string   `xml:"Id"`
	Title   string   `xml:"Title"`
}

type PersonConfidential struct {
	XMLName xml.Name        `xml:"PersonConfidential"`
	Row     ConfidentialRow `xml:"Row"`
}

type ConfidentialRow struct {
	XMLName      xml.Name `xml:"Row"`
	What         string   `xml:"What"`
	Value        string   `xml:"Value"`
	IsAllowedCST string   `xml:"IsAllowedCST"`
	IsAllowedADS string   `xml:"IsAllowedADS"`
	IsAllowedTB  string   `xml:"IsAllowedTB"`
}
