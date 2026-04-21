package getaccinfo

import "encoding/xml"

// Envelope - корневой элемент SOAP конверта
type Envelope struct {
	XMLName xml.Name `xml:"s:Envelope"`
	XmlnsS  string   `xml:"xmlns:s,attr"`
	XmlnsM1 string   `xml:"xmlns:m1,attr"`
	XmlnsM0 string   `xml:"xmlns:m0,attr"`
	Body    RespBody `xml:"s:Body"`
}

type RespBody struct {
	GetAcctInfoRp GetAcctInfoRp `xml:"m1:GetAcctInfoRp"`
}

type GetAcctInfoRp struct {
	Response Response `xml:"m1:Response"`
}

type Response struct {
	Echo         string `xml:"Echo,attr"`
	Product      string `xml:"Product,attr"`
	ResponseAttr string `xml:"Response,attr"`
	TranId       string `xml:"TranId,attr"`
	Ver          string `xml:"Ver,attr"`

	Avail                 string `xml:"m0:Avail"`
	Bonus                 string `xml:"m0:Bonus"`
	Cards                 Rows   `xml:"m0:Cards"`
	CreditHold            string `xml:"m0:CreditHold"`
	Currency              string `xml:"m0:Currency"`
	DebitHold             string `xml:"m0:DebitHold"`
	DropTmpOverOnRefresh  string `xml:"m0:DropTmpOverOnRefresh"`
	ExtendedAccountNumber string `xml:"m0:ExtendedAccountNumber"`
	FoundAccount          string `xml:"m0:FoundAccount"`
	LastDepAmount         string `xml:"m0:LastDepAmount"`
	LastDepTime           string `xml:"m0:LastDepTime"`     //2006-01-02T15:04:05
	LastRefreshTime       string `xml:"m0:LastRefreshTime"` //2006-01-02T15:04:05
	LastTranId            string `xml:"m0:LastTranId"`
	LastWdlAmount         string `xml:"m0:LastWdlAmount"`
	LastWdlTime           string `xml:"m0:LastWdlTime"`
	Ledger                string `xml:"m0:Ledger"`
	MaskBalances          string `xml:"m0:MaskBalances"`
	PermissibleExcessType string `xml:"m0:PermissibleExcessType"`
	PersonExtId           string `xml:"m0:PersonExtId"`
	PersonFIO             string `xml:"m0:PersonFIO"`
	PersonId              string `xml:"m0:PersonId"`
	Remain                string `xml:"m0:Remain"`
	Status                string `xml:"m0:Status"`
	TmpOverdraft          string `xml:"m0:TmpOverdraft"`
	Type                  string `xml:"m0:Type"`
}

type Rows struct {
	Rows []CardRow `xml:"m0:Row"`
}

type CardRow struct {
	PAN    string `xml:"m0:PAN"`
	MBR    string `xml:"m0:MBR"`
	Status string `xml:"m0:Status"`
	Type   string `xml:"m0:Type"`
}

// <s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:m1="http://schemas.compassplus.com/two/1.0/fimi.xsd" xmlns:m0="http://schemas.compassplus.com/two/1.0/fimi_types.xsd">
//    <s:Body>
//       <m1:GetAcctInfoRp>
//          <m1:Response Echo="" Encoding="" Product="FIMI" Response="1" Ver="16.37">
//             <m0:Avail>1</m0:Avail>
//             <m0:Bonus>0</m0:Bonus>
//             <m0:Cards>
//                <m0:Row>
//                   <m0:PAN>9762******4101</m0:PAN>
//                   <m0:MBR>0</m0:MBR>
//                   <m0:Status>1</m0:Status>
//                   <m0:Type>1</m0:Type>
//                </m0:Row>
//             </m0:Cards>
//             <m0:CreditHold>0</m0:CreditHold>
//             <m0:Currency>972</m0:Currency>
//             <m0:DebitHold>0</m0:DebitHold>
//             <m0:DropTmpOverOnRefresh>0</m0:DropTmpOverOnRefresh>
//             <m0:ExtendedAccountNumber>20216972*****048</m0:ExtendedAccountNumber>
//             <m0:FoundAccount>20216972*****048</m0:FoundAccount>
//             <m0:LastDepAmount>3041</m0:LastDepAmount>
//             <m0:LastDepTime>2024-11-11T20:49:59</m0:LastDepTime>
//             <m0:LastRefreshTime>2026-02-13T08:18:45</m0:LastRefreshTime>
//             <m0:LastTranId>92***2833</m0:LastTranId>
//             <m0:LastWdlAmount>40</m0:LastWdlAmount>
//             <m0:LastWdlTime>2025-03-27T14:47:20</m0:LastWdlTime>
//             <m0:Ledger>1</m0:Ledger>
//             <m0:MaskBalances>0</m0:MaskBalances>
//             <m0:PermissibleExcessType>-1</m0:PermissibleExcessType>
//             <m0:PersonExtId>2015838</m0:PersonExtId>
//             <m0:PersonFIO>992****60 ХУДОЯРОВА МУХБИРА МУХАМАДОВНА</m0:PersonFIO>
//             <m0:PersonId>2237824</m0:PersonId>
//             <m0:Remain>1</m0:Remain>
//             <m0:Status>3</m0:Status>
//             <m0:TmpOverdraft>0</m0:TmpOverdraft>
//             <m0:Type>1</m0:Type>
//          </m1:Response>
//       </m1:GetAcctInfoRp>
//    </s:Body>
// </s:Envelope>
