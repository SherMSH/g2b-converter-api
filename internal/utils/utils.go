package utils

import (
	"converterapi/pkg/logger"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func Init() {
	D8HeadersMap["Pragma"] = "no-cache"
	D8HeadersMap["Cashe-Control"] = "no-cache"
	D8HeadersMap["Content-Type"] = "application/json"
}

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

func SendRequest(method, uri string, jsonBody []byte, headers map[string]string) (data []byte, status int, err error) {

	req, err := http.NewRequest(method, uri, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, 500, fmt.Errorf("Request err: %v", err)
	}

	logger.Infof("Sending G2b req: %s %s", method, uri)

	client := http.Client{
		Timeout: 90 * time.Second,
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Close = true

	resps, err := client.Do(req)
	if err != nil {
		return nil, 500, fmt.Errorf("Request sending error: %v", err)
	}
	defer resps.Body.Close()

	status = resps.StatusCode
	body, err := io.ReadAll(resps.Body)
	if err != nil {
		return nil, 500, fmt.Errorf("Response body reading error: %v", err)
	}

	return body, status, nil
}
