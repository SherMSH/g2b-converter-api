package service

import (
	"bytes"
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	d8procweb "converterapi/pkg/d8-proc-web"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
	"io"
)

func GetAcctInfoG2b(accNum string) (foundAcc *d8procweb.AccountData, err error) {
	var (
		accs []d8procweb.AccountData
	)
	resp, err := d8procweb.Signin()
	if err != nil {
		return nil, fmt.Errorf("d8procweb signin error %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("d8procweb signin staus: %v", resp.StatusCode)
	}
	defer d8procweb.Signout()

	reqBody := d8procweb.RequestBody{
		Data:     map[string]interface{}{},
		Ordering: []string{},
		Filters: []d8procweb.Filter{
			{
				Column: "accnum",
				Values: []string{accNum},
			},
		},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marchaling reuest body: %v", err)
	}
	resp, err = d8procweb.Client.Post(config.Config.Processing.Address+"/api/miss/v1/getAccountList", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("Reuest Error %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %v", err)
	}

	logger.Debugf("getAccountList Статус: %s", resp.Status)
	logger.Debugf("getAccountList Ответ: %s", string(body))

	getAccsResp := d8corp.CommonResp{}
	err = json.Unmarshal(body, &getAccsResp)
	if err != nil {
		return nil, fmt.Errorf("getAccountList resp unmarshaling error: %v", err)
	}
	if getAccsResp.Status.RspCode != "00" || resp.StatusCode != 200 {
		return nil, fmt.Errorf("getAccountList resp status error: %v", err)
	}
	err = json.Unmarshal(getAccsResp.Data, &accs)
	if err != nil {
		return nil, fmt.Errorf("Accounts data unmarshal error: %v", err)
	}
	if len(accs) == 0 {
		return nil, fmt.Errorf("Account %v not found", accNum)
	}

	reqBody = d8procweb.RequestBody{
		Data: d8procweb.AccountData{
			ID: accs[0].ID,
		},
	}

	jsonData, err = json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request body error: %v", err)
	}
	resp, err = d8procweb.Client.Post(config.Config.Processing.Address+"/api/miss/v1/getAccount", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("getAccount Reqest Error %v", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %v", err)
	}

	logger.Debugf("getAccount Статус: %s", resp.Status)
	logger.Debugf("getAccount Ответ: %s", string(body))

	getAccResp := d8corp.CommonResp{}
	err = json.Unmarshal(body, &getAccResp)
	if err != nil {
		return nil, fmt.Errorf("getAccount resp unmarshaling error: %v", err)
	}
	if getAccResp.Status.RspCode != "00" || resp.StatusCode != 200 {
		return nil, fmt.Errorf("getAccount resp status error: %v", err)
	}

	err = json.Unmarshal(getAccResp.Data, &foundAcc)
	if err != nil {
		return nil, fmt.Errorf("Accounts data unmarshal error: %v", err)
	}
	if foundAcc.ID == 0 {
		return nil, fmt.Errorf("Account %v not found", accNum)
	}
	foundAcc.Name = accs[0].Name
	foundAcc.Custcode = accs[0].Custcode
	foundAcc.CompanyRegnum = accs[0].CompanyRegnum
	foundAcc.CompanyName = accs[0].CompanyName

	return foundAcc, nil
}
