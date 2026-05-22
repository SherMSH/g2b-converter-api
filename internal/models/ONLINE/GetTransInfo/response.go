package gettransinfo

import "encoding/xml"

type Envelope struct {
	XMLName xml.Name `xml:"s:Envelope"`
	XmlnsS  string   `xml:"xmlns:s,attr"`
	XmlnsM1 string   `xml:"xmlns:m1,attr"`
	XmlnsM0 string   `xml:"xmlns:m0,attr"`
	Body    RespBody `xml:"s:Body"`
}

type RespBody struct {
	GetTransInfoRp GetTransInfoRp `xml:"m1:GetTransInfoRp"`
}

type GetTransInfoRp struct {
	Response Response `xml:"m1:Response"`
}

type Response struct {
	Echo         string `xml:"Echo,attr"`
	Product      string `xml:"Product,attr"`
	ResponseAttr string `xml:"Response,attr"`
	Ver          string `xml:"Ver,attr"`

	MaskBalances string   `xml:"m0:MaskBalances"`
	TranList     TranList `xml:"m0:TranList"`
}

type TranList struct {
	Rows []TranListRow `xml:"m0:Row"`
}

type TranListRow struct {
	Id                   string `xml:"m0:Id"`
	Type                 string `xml:"m0:Type"`
	Time                 string `xml:"m0:Time"`
	Phase                string `xml:"m0:Phase"`
	TermClass            string `xml:"m0:TermClass"`
	TermName             string `xml:"m0:TermName"`
	TermDate             string `xml:"m0:TermDate"`
	TranCode             string `xml:"m0:TranCode"`
	DraftCapture         string `xml:"m0:DraftCapture"`
	FromAcct             string `xml:"m0:FromAcct"`
	Amount               string `xml:"m0:Amount"`
	Fee                  string `xml:"m0:Fee"`
	Issuer               string `xml:"m0:IssuerFee"`
	Currency             string `xml:"m0:Currency"`
	PAN                  string `xml:"m0:PAN"`
	Card                 string `xml:"m0:CardMember"`
	Resp                 string `xml:"m0:RespCode"`
	Retain               string `xml:"m0:RetainCard"`
	Approval             string `xml:"m0:ApprovalCode"`
	LedgerBalance        string `xml:"m0:LedgerBalance"`
	AvailBalance         string `xml:"m0:AvailBalance"`
	BalanceCurrencyAcct  string `xml:"m0:BalanceCurrencyAcct"`
	CurrencyAcct         string `xml:"m0:CurrencyAcct"`
	AmountAcct           string `xml:"m0:AmountAcct"`
	Exchange             string `xml:"m0:ExchangeRateAcct"`
	RevRequestId         string `xml:"m0:RevRequestId"`
	Error                string `xml:"m0:Error"`
	OrigType             string `xml:"m0:OrigType"`
	TermFIName           string `xml:"m0:TermFIName"`
	TermInstID           string `xml:"m0:TermInstID"`
	TermRetailerName     string `xml:"m0:TermRetailerName"`
	TermSIC              string `xml:"m0:TermSIC"`
	TermSICName          string `xml:"m0:TermSICName"`
	TermCountry          string `xml:"m0:TermCountry"`
	TermCountryName      string `xml:"m0:TermCountryName"`
	TermCity             string `xml:"m0:TermCity"`
	TermOwner            string `xml:"m0:TermOwner"`
	Track2               string `xml:"m0:Track2"`
	Auth                 string `xml:"m0:AuthFIName"`
	RevActualAmount      string `xml:"m0:RevActualAmount"`
	TranNumber           string `xml:"m0:TranNumber"`
	POSCondition         string `xml:"m0:POSCondition"`
	POSEntry             string `xml:"m0:POSEntryMode"`
	From                 string `xml:"m0:FromAcctType"`
	To                   string `xml:"m0:ToAcctType"`
	CNSent               string `xml:"m0:CNSent"`
	Overdraft            string `xml:"m0:OverdraftLimit"`
	Tmp                  string `xml:"m0:TmpOverdraft"`
	Prev                 string `xml:"m0:PrevTran"`
	Orig                 string `xml:"m0:OrigTime"`
	DebitHold            string `xml:"m0:DebitHold"`
	CreditHold           string `xml:"m0:CreditHold"`
	Bonus                string `xml:"m0:Bonus"`
	LedgerBalanceBefore  string `xml:"m0:LedgerBalanceBefore"`
	AvailBalanceBefore   string `xml:"m0:AvailBalanceBefore"`
	DebitHoldBefore      string `xml:"m0:DebitHoldBefore"`
	CreditHoldBefore     string `xml:"m0:CreditHoldBefore"`
	BonusBefore          string `xml:"m0:BonusBefore"`
	Reason               string `xml:"m0:Reason"`
	ICC_IssuerScript1    string `xml:"m0:ICC_IssuerScript1"`
	ICC_IssuerScript2    string `xml:"m0:ICC_IssuerScript2"`
	ICC_TranType         string `xml:"m0:ICC_TranType"`
	Clerk                string `xml:"m0:Clerk"`
	PINEntry             string `xml:"m0:PINEntry"`
	Host                 string `xml:"m0:Host"`
	CNSId                string `xml:"m0:CNSId"`
	LedgerBalance2       string `xml:"m0:LedgerBalance2"`
	AvailBalance2        string `xml:"m0:AvailBalance2"`
	LedgerBalance2Before string `xml:"m0:LedgerBalance2Before"`
	AvailBalance2Before  string `xml:"m0:AvailBalance2Before"`
	ReversalAllowed      string `xml:"m0:ReversalAllowed"`
	ReceiptFlag          string `xml:"m0:ReceiptFlag"`
	IsContainCAVV        string `xml:"m0:IsContainCAVV"`
	TermInstCountry      string `xml:"m0:TermInstCountry"`
	TermInstCountryName  string `xml:"m0:TermInstCountryName"`
	TranTime             string `xml:"m0:TranTime"`
	RevActualAmountAcct  string `xml:"m0:RevActualAmountAcct"`
}
