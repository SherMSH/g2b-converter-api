package getcardinfo

import "encoding/xml"

type Envelope struct {
	XMLName xml.Name `xml:"s:Envelope"`
	XmlnsS  string   `xml:"xmlns:s,attr"`
	XmlnsM1 string   `xml:"xmlns:m1,attr"`
	XmlnsM0 string   `xml:"xmlns:m0,attr"`
	Body    RespBody `xml:"s:Body"`
}

type RespBody struct {
	GetCardInfoRp GetCardInfoRp `xml:"m1:GetCardInfoRp"`
}

type GetCardInfoRp struct {
	Response Response `xml:"m1:Response"`
}

type Response struct {
	Echo         string `xml:"Echo,attr"`
	Product      string `xml:"Product,attr"`
	ResponseAttr string `xml:"Response,attr"`
	TranId       string `xml:"TranId,attr"`
	Ver          string `xml:"Ver,attr"`

	Accounts                Accounts           `xml:"m0:Accounts"`
	Acct2CardAttachType     string             `xml:"m0:Acct2CardAttachType"`
	CNSDisabled             string             `xml:"m0:CNSDisabled"`
	CardAllowedEMVScript    string             `xml:"m0:CardAllowedEMVScript"`
	CardProfiles            CardProfiles       `xml:"m0:CardProfiles"`
	ContactlessStatus       string             `xml:"m0:ContactlessStatus"`
	ECNeedCAPAuth           string             `xml:"m0:ECNeedCAPAuth"`
	ECNeedDynPwdAuth        string             `xml:"m0:ECNeedDynPwdAuth"`
	ECNeedStaticAuth        string             `xml:"m0:ECNeedStaticAuth"`
	ECNeedTokenAuth         string             `xml:"m0:ECNeedTokenAuth"`
	ECStatus                string             `xml:"m0:ECStatus"`
	ECUseCardSettingsAuth   string             `xml:"m0:ECUseCardSettingsAuth"`
	ECUseDecoupledAuth      string             `xml:"m0:ECUseDecoupledAuth"`
	EMVOptionsCheckDisabled string             `xml:"m0:EMVOptionsCheckDisabled"`
	ExpDate                 string             `xml:"m0:ExpDate"`
	FoundMBR                string             `xml:"m0:FoundMBR"`
	FoundPAN                string             `xml:"m0:FoundPAN"`
	IB_Registered           string             `xml:"m0:IB_Registered"`
	InstName                string             `xml:"m0:InstName"`
	IssueTechnology         string             `xml:"m0:IssueTechnology"`
	LastATMUsed             string             `xml:"m0:LastATMUsed"`
	LastChangeStatusTime    string             `xml:"m0:LastChangeStatusTime"`
	LastPOSUsed             string             `xml:"m0:LastPOSUsed"`
	LastPVVChangeTime       string             `xml:"m0:LastPVVChangeTime"`
	LastRefreshTime         string             `xml:"m0:LastRefreshTime"`
	LastTranId              string             `xml:"m0:LastTranId"`
	LastTranTime            string             `xml:"m0:LastTranTime"`
	MaskBalances            string             `xml:"m0:MaskBalances"`
	MaskPVV                 string             `xml:"m0:MaskPVV"`
	NameOnCard              string             `xml:"m0:NameOnCard"`
	PINVerifyType           string             `xml:"m0:PINVerifyType"`
	PVV                     string             `xml:"m0:PVV"`
	PasswordFlag            string             `xml:"m0:PasswordFlag"`
	PersonConfidential      PersonConfidential `xml:"m0:PersonConfidential"`
	PersonExtId             string             `xml:"m0:PersonExtId"`
	PersonFIO               string             `xml:"m0:PersonFIO"`
	PersonId                string             `xml:"m0:PersonId"`
	PersonVIP               string             `xml:"m0:PersonVIP"`
	RequiredPasswordVersion string             `xml:"m0:RequiredPasswordVersion"`
	RiskControlDisabled     string             `xml:"m0:RiskControlDisabled"`
	RiskLevel               string             `xml:"m0:RiskLevel"`
	Status                  string             `xml:"m0:Status"`
	TmpECStatus             string             `xml:"m0:TmpECStatus"`
	Type                    string             `xml:"m0:Type"`
	UseUdCVV2               string             `xml:"m0:UseUdCVV2"`
}

type Accounts struct {
	XMLName xml.Name     `xml:"m0:Accounts"`
	Row     []AccountRow `xml:"m0:Row"`
}

type AccountRow struct {
	XMLName       xml.Name `xml:"m0:Row"`
	AcctNo        string   `xml:"m0:AcctNo"`
	Status        string   `xml:"m0:Status"`
	LedgerBalance string   `xml:"m0:LedgerBalance"`
	AvailBalance  string   `xml:"m0:AvailBalance"`
	Currency      string   `xml:"m0:Currency"`
	Type          string   `xml:"m0:Type"`
	AccountStatus string   `xml:"m0:AccountStatus"`
}

type CardProfiles struct {
	XMLName xml.Name       `xml:"m0:CardProfiles"`
	Row     CardProfileRow `xml:"m0:Row"`
}

type CardProfileRow struct {
	XMLName xml.Name `xml:"m0:Row"`
	Id      string   `xml:"m0:Id"`
	Title   string   `xml:"m0:Title"`
}

type PersonConfidential struct {
	Row ConfidentialRow `xml:"m0:Row"`
}

type ConfidentialRow struct {
	XMLName      xml.Name `xml:"m0:Row"`
	What         string   `xml:"m0:What"`
	Value        string   `xml:"m0:Value"`
	IsAllowedCST string   `xml:"m0:IsAllowedCST"`
	IsAllowedADS string   `xml:"m0:IsAllowedADS"`
	IsAllowedTB  string   `xml:"m0:IsAllowedTB"`
}
