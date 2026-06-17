package service

import (
	"bytes"
	"converterapi/internal/config"
	d8procweb "converterapi/pkg/d8-proc-web"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
	"io"
)

type Filter struct {
	Column string   `json:"column"`
	Values []string `json:"values"`
}

type RequestBody struct {
	Data     interface{} `json:"data"`
	Ordering []string    `json:"ordering"`
	Filters  []Filter    `json:"filters"`
}

func GetAcctInfoG2b(accNum string) error {
	resp, err := d8procweb.Signin()
	if err != nil {
		return fmt.Errorf("d8procweb signin error %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("d8procweb signin staus: %v", resp.StatusCode)
	}
	defer d8procweb.Signout()

	reqBody := RequestBody{
		Data:     map[string]interface{}{},
		Ordering: []string{},
		Filters: []Filter{
			{
				Column: "accnum",
				Values: []string{accNum},
			},
		},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("error marchaling reuest body: %v", err)
	}
	resp, err = d8procweb.Client.Post(config.Config.Processing.Address+"/api/miss/v1/getAccountList", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Reuest Error %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading response: %v", err)
	}

	logger.Debugf("Статус: %s", resp.Status)
	logger.Debugf("Ответ: %s", string(body))

	return nil
}
