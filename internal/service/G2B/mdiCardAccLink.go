package service

import (
	"converterapi/internal/config"
	"converterapi/internal/models"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
)

func AddCardAcctLinkG2b(input models.MDIface) (resp interface{}, err error) {
	var recDetails d8corp.MdiFile
	recNums := utils.NewSequence()

	for _, v := range input.GetRecords() {
		record := d8corp.MdiRecordDetails{
			IssRectype:         "CRDACC",
			IssRecaction:       "ADD",
			IssRecnum:          recNums.NextVal(),
			IssCompanyRegnr:    "ARVD",
			IssCompanyRegnrAcc: "ARVD",
			KlLKeyClr:          v.PAN,
			DbAccountCurrcode:  v.CurrencyNo,
			DbAccountAccnum:    v.Account,
		}
		jsonRec, err := json.Marshal(record)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b ADD CRDACC req marshaling record err: %v", err)
			return nil, err
		}
		recDetails.MdiRecords = append(recDetails.MdiRecords, jsonRec)
	}

	reqJSON, err := json.Marshal(recDetails)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADD CRDACC req marshaling err: %v", err)
		return nil, err
	}
	logger.Infof("json ADD CRDACC: %v", string(reqJSON))
	mdiFile := d8corp.MdiFile{
		MdiRecords: []json.RawMessage{
			reqJSON,
		},
	}
	mdiDataJSON, err := json.MarshalIndent(mdiFile, "", "  ")
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADD CRDACC req marshaling err: %v", err)
		return nil, err
	}

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/mdi", mdiDataJSON, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADD CRDACC request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b ADD CRDACC resp status: %v, body: %v", status, string(data))
	return data, nil
}

func UpdateCardAcctLinkG2b(input models.MDIface) (resp interface{}, err error) {
	var recDetails d8corp.MdiFile
	recNums := utils.NewSequence()

	for _, v := range input.GetRecords() {
		record := d8corp.MdiRecordDetails{
			IssRectype:         "CRDACC",
			IssRecaction:       "UPDATE",
			IssRecnum:          recNums.NextVal(),
			IssCompanyRegnr:    "ARVD",
			IssCompanyRegnrAcc: "ARVD",
			KlLKeyClr:          v.PAN,
			DbAccountCurrcode:  v.CurrencyNo,
			DbAccountAccnum:    v.Account,
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

func DeleteCardAcctLinkG2b(input models.MDIface) (resp interface{}, err error) {
	var recDetails d8corp.MdiFile
	recNums := utils.NewSequence()

	for _, v := range input.GetRecords() {
		record := d8corp.MdiRecordDetails{
			IssRectype:         "CRDACC",
			IssRecaction:       "DELETE",
			IssRecnum:          recNums.NextVal(),
			IssCompanyRegnr:    "ARVD",
			IssCompanyRegnrAcc: "ARVD",
			KlLKeyClr:          v.PAN,
			DbAccountCurrcode:  v.CurrencyNo,
			DbAccountAccnum:    v.Account,
		}
		jsonRec, err := json.Marshal(record)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b DELETE CRDACC req marshaling record err: %v", err)
			return nil, err
		}
		recDetails.MdiRecords = append(recDetails.MdiRecords, jsonRec)
	}

	reqJSON, err := json.Marshal(recDetails)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b DELETE CRDACC req marshaling err: %v", err)
		return nil, err
	}
	logger.Infof("json DELETE CRDACC: %v", string(reqJSON))
	mdiFile := d8corp.MdiFile{
		MdiRecords: []json.RawMessage{
			reqJSON,
		},
	}
	mdiDataJSON, err := json.MarshalIndent(mdiFile, "", "  ")
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b DELETE CRDACC req marshaling err: %v", err)
		return nil, err
	}

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/mdi", mdiDataJSON, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b DELETE CRDACC request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b DELETE CRDACC resp status: %v, body: %v", status, string(data))
	return data, nil
}
