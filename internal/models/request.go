package models

import (
	models "converterapi/internal/models/OFFLINE"
	"converterapi/internal/utils"
	"encoding/xml"
	"reflect"
)

// SoapEnvelope соответствует корневому элементу <Envelope>
type SoapEnvelope struct {
	XMLName xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	Body    XmlBody  `xml:"Body"`
}

type XmlBody struct {
	XMLData []byte `xml:",innerxml" json:"body"`
}

type SoapBody interface {
	GetBodyType() reflect.Type
}

type TrnInputIface interface {
	GetTxnType() utils.TxnType
	GetPan() string
	GetMBR() string
	GetExpDate() string
	GetAccount() string
	GetAmount() float64
	GetCurrency() string
	GetRecipientAcc() string
	GetTerminal() string
	GetAcceptorID() string
}

type MDIface interface {
	GetRecords() []models.MRecord
	GetRecordsCount() int
}
