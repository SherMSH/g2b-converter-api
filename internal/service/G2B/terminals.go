package service

import (
	d8procweb "converterapi/pkg/d8-proc-web"
	"encoding/json"
)

func GetTerminalByIdG2b(tid string) (foundTerminals []d8procweb.TerminalsData, err error) {
	path := "/api/pos/v1/getTerminalList"
	filters := []d8procweb.Filter{
		{
			Column: "terminalid",
			Values: []string{tid},
		},
	}
	d8procweb.Signin()
	defer d8procweb.Signout()
	resp, err := d8procweb.Request(path, filters)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &foundTerminals)
	if err != nil {
		return nil, err
	}
	return
}
