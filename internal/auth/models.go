package auth

import (
	"encoding/xml"
	"strings"
)

// Структура только для получения Clerk и Password
type AuthRequest struct {
	XMLName  xml.Name `xml:"Request"`
	Ver      string   `xml:"Ver,attr"`
	Product  string   `xml:"Product,attr"`
	Echo     string   `xml:"Echo,attr"`
	Session  string   `xml:"Session,attr"`
	Clerk    string   `xml:"Clerk,attr"`
	Password string   `xml:"Password,attr"`
}

// Используем xml:",any" для пропуска промежуточных элементов
type SOAPEnvelopeSimple struct {
	XMLName xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	Body    struct {
		XMLName xml.Name `xml:"Body"`
		Any     struct {
			XMLName xml.Name
			Request AuthRequest `xml:",any"` // Игнорируем промежуточные теги
		} `xml:",any"`
	} `xml:"Body"`
}

func ParseAuth(xmlData string) (*AuthRequest, error) {
	var envelope SOAPEnvelopeSimple

	decoder := xml.NewDecoder(strings.NewReader(xmlData))
	decoder.Strict = false // Отключаем строгий режим для игнорирования namespace

	if err := decoder.Decode(&envelope); err != nil {
		return nil, err
	}

	return &envelope.Body.Any.Request, nil
}

var ClerksMap map[string]string = map[string]string{
	"davronjon.boltaboev@arvand.tj": "4799726cbd1037e3a3d1f767ca1a4d297c9fb354664e02ad4a1573f4ebede80b",
}
