package d8procweb

import (
	"bytes"
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
	"io"
)

func D8procwebRequest(path string, filters []Filter) (respData json.RawMessage, err error) {
	reqBody := RequestBody{
		Filters: filters,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marchaling request body: %v", err)
	}
	if Client == nil {
		Init()
	}
	resp, err := Client.Post(config.Config.Processing.Address+path, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("%s Request Error %v", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %v", err)
	}

	logger.Debugf("[SERVICE] %s Статус: %s", path, resp.Status)
	logger.Debugf("[SERVICE] %s Ответ: %s", path, string(body))

	commonResp := d8corp.CommonResp{}
	err = json.Unmarshal(body, &commonResp)
	if err != nil {
		return nil, fmt.Errorf("%s resp unmarshaling error: %v", path, err)
	}
	if commonResp.Status.RspCode != "00" || resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s resp status error: %v", path, err)
	}

	return commonResp.Data, nil
}
