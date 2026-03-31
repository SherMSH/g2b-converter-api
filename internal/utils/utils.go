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

	D8TxHeadersMap["Pragma"] = "no-cache"
	D8TxHeadersMap["Cashe-Control"] = "no-cache"
	D8TxHeadersMap["Content-Type"] = "application/json"
	D8TxHeadersMap["Connection"] = "keep-alive"
	D8TxHeadersMap["X-Ssl-Client-Serial"] = "0B6F8A39AAC45D27F61F7E1A2D2F94CC"
	D8TxHeadersMap["X-Ssl-Client-I-Dn"] = "CN=Test Company CA, O=Test Company, C=US"
	D8TxHeadersMap["X-Ssl-Client-S-Dn"] = "CN=John Doe, O=Test Company, C=US"
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
		// logger.Infof("Hey: %v value: %v", k, v)
		req.Header.Set(k, v)
	}
	// req.Close = true
	logger.Infof("url: %v, headers: %v, body: %v", req.URL, req.Header, string(jsonBody))

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

func GetExpFromTrack(trck2 string) string {
	if !strings.Contains(trck2, "=") {
		return ""
	}
	data := strings.Split(trck2, "=")
	if len(data) < 2 {
		return ""
	}
	if len(data[2]) <= 4 {
		return data[2]
	}
	return data[2][:4]
}
