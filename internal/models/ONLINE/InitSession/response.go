package initsession

import "encoding/xml"

type Envelope struct {
	XMLName xml.Name `xml:"s:Envelope"`
	XmlnsS  string   `xml:"xmlns:s,attr"`
	XmlnsM1 string   `xml:"xmlns:m1,attr"`
	XmlnsM0 string   `xml:"xmlns:m0,attr"`
	Body    RespBody `xml:"s:Body"`
}

type RespBody struct {
	InitSessionRp InitSessionRp `xml:"m1:InitSessionRp"`
}

type InitSessionRp struct {
	Response Response `xml:"m1:Response"`
}

type Response struct {
	Echo          string `xml:"Echo,attr"`
	Product       string `xml:"Product,attr"`
	ResponseAttr  string `xml:"Response,attr"`
	Ver           string `xml:"Ver,attr"`
	NextChallenge string `xml:"NextChallenge,attr"`

	CNSChannelList  CNSChannelList `xml:"m0:CNSChannelList"`
	CNSSchemeList   CNSSchemeList  `xml:"m0:CNSSchemeList"`
	Id              string         `xml:"m0:Id"`
	NeedCAPAuth     string         `xml:"m0:NeedCAPAuth"`
	PasswordVersion string         `xml:"m0:PasswordVersion"`
}

type CNSChannelList struct {
	Rows []CNSChannelRow `xml:"m0:Row"`
}

type CNSSchemeList struct {
	Rows []CNSSchemeRow `xml:"m0:Row"`
}

type CNSChannelRow struct {
	Name        string `xml:"m0:Name"`
	Title       string `xml:"m0:Title"`
	UsedForPush string `xml:"m0:UsedForPush"`
}
type CNSSchemeRow struct {
	Name  string `xml:"m0:Name"`
	Title string `xml:"m0:Title"`
}
