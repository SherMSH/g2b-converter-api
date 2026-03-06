package gettransinfo

import "reflect"

type Body struct {
	SoapRq SoapRq `xml:"GetTransInfoRq" json:"GetTransInfoRq"`
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

	FromTime    string `xml:"FromTime" json:"from_time"`
	ToTime      string `xml:"ToTime" json:"to_time"`
	MBR         string `xml:"MBR" json:"mbr"`
	OrigType    string `xml:"OrigType" json:"orig_type"`
	TermName    string `xml:"TermName" json:"term_name"`
	TranNumber  string `xml:"TranNumber" json:"tran_number"`
	AcqInstName string `xml:"AcqInstName" json:"acq_inst_name"`

	Count          string `xml:"Count" json:"count"`
	Id             string `xml:"Id" json:"id"`
	PAN            string `xml:"PAN" json:"pan"`
	PartitionName  string `xml:"PartitionName" json:"partition_name"`
	TermId         string `xml:"TermId" json:"term_id"`
	IssInstName    string `xml:"IssInstName" json:"is_inst_name"`
	CardUID        string `xml:"CardUID" json:"card_uid"`
	ResponseFields string `xml:"ResponseFields" json:"response_fields"`
	BusinessDay    string `xml:"BusinessDay" json:"business_day"`
	NewerTran      string `xml:"NewerTran" json:"newer_tran"`
	ExtRRN         string `xml:"ExtRRN" json:"ext_rrn"`

	Type     Rows `xml:"Type" json:"type"`
	TranCode Rows `xml:"TranCode" json:"tran_code"`

	ArrayTermId      Rows `xml:"ArrayTermId" json:"array_term_id"`
	ArrayAcqInstName Rows `xml:"ArrayAcqInstName" json:"array_acq_inst_name"`
	ArrayIssInstName Rows `xml:"ArrayIssInstName" json:"array_iss_inst_name"`
}

type Rows struct {
	Row []Row `xml:"Row" json:"rows"`
}

type Row struct {
	Value string `xml:"Value" json:"value"`
}
