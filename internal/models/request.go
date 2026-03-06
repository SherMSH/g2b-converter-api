package models

import (
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
