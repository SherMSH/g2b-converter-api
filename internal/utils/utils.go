package utils

import (
	"encoding/xml"
	"strings"
)

func GetRqType(xmlData string) RqBodyType {
	decoder := xml.NewDecoder(strings.NewReader(xmlData))
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			// Проверяем локальное имя элемента
			for _, v := range BodyTypes {
				if se.Name.Local == string(v) {
					return v
				}
			}
		}
	}
	return Unknown
}
