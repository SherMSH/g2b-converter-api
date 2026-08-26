package reversetransaction

import "reflect"

type Body struct {
	SoapRq SoapRq `xml:"ReverseTransactionRq" json:"ReverseTransactionRq"`
}

func (sb Body) GetBodyType() reflect.Type {
	return reflect.TypeOf(sb)
}

func (sb *Body) Call() (*Envelope, error) {
	return Svc(sb)
}

// SoapRq соответствует элементу fimi:ReverseTransactionRq
type SoapRq struct {
	Req Request `xml:"Request" json:"Request"`
}

// Request соответствует элементу fimi:Request.
//
// Партнёр передаёт только идентификатор отменяемой операции; сумма, валюта и
// ссылка на транзакцию в процессинге достаются по нему из getTransactionDetails.
type Request struct {
	Ver               string `xml:"Ver,attr" json:"ver"`
	Product           string `xml:"Product,attr" json:"product"`
	Echo              string `xml:"Echo,attr" json:"echo"`
	Session           string `xml:"Session,attr" json:"session"`
	Clerk             string `xml:"Clerk,attr" json:"clerk"`
	Password          string `xml:"Password,attr" json:"password"`
	TransactionNumber string `xml:"TransactionNumber,attr" json:"transaction_number"`

	TranNumber string `xml:"TranNumber" json:"tran_number"`
	PAN        string `xml:"PAN" json:"pan"`
	MBR        string `xml:"MBR" json:"mbr"`
	Track2     string `xml:"Track2" json:"track2"`

	// Id         string `xml:"Id" json:"id"`
	// TermName   string `xml:"TermName" json:"term_name"`
	// SettleDate string `xml:"SettleDate" json:"settle_date"`
	Amount     string `xml:"Amount" json:"amount"`
	ReasonCode string `xml:"ReasonCode" json:"reason_code"`
}
