package service

import (
	"converterapi/internal/config"
	"converterapi/internal/models"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"strconv"
)

func GetCustomerInfoG2b() error {
	return nil
}

func UpdateCustomerG2b(input models.MDIface) (resp interface{}, err error) {
	var recDetails d8corp.MdiFile
	recNums := utils.NewSequence()

	for _, v := range input.GetRecords() {
		homeCountry, err := strconv.Atoi(v.CntrLive)
		if err != nil {
			homeCountry = 972
			err = nil
		}

		record := d8corp.MdiRecordDetails{
			IssRectype:            "CUSTOMER",
			IssRecaction:          "UPDATE",
			IssRecnum:             recNums.NextVal(),
			IssCompanyRegnr:       "ARVD",
			DbCustomerCustcode:    v.Id,
			DbCustomerFirstName:   v.FirstName,
			DbCustomerLastName:    v.LastName,
			DbCardFathername:      v.MiddleName,
			DbCustomerDateBirth:   v.Birthday,
			DbCustomerHomeCountry: homeCountry,
		}
		jsonRec, err := json.Marshal(record)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b UPDATE CRDACC req marshaling record err: %v", err)
			return nil, err
		}
		recDetails.MdiRecords = append(recDetails.MdiRecords, jsonRec)
	}

	reqJSON, err := json.Marshal(recDetails)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b UPDATE CRDACC req marshaling err: %v", err)
		return nil, err
	}
	logger.Infof("json UPDATE CRDACC: %v", string(reqJSON))
	mdiFile := d8corp.MdiFile{
		MdiRecords: []json.RawMessage{
			reqJSON,
		},
	}
	mdiDataJSON, err := json.MarshalIndent(mdiFile, "", "  ")
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b UPDATE CRDACC req marshaling err: %v", err)
		return nil, err
	}

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/mdi", mdiDataJSON, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b UPDATE CRDACC request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b UPDATE CRDACC resp status: %v, body: %v", status, string(data))
	return data, nil

}
