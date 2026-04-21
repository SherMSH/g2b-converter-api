package service

import (
	"converterapi/internal/config"
	"converterapi/internal/models"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
)

func AddCardNotificationG2b(input models.MDIface) (result *d8corp.CommonResp, err error) {
	var (
		recDetails d8corp.MdiFile
		resp       d8corp.CommonResp
		mdiData    d8corp.MdiData
	)
	recNums := utils.NewSequence()

	for _, v := range input.GetRecords() {
		record := d8corp.MdiRecordDetails{
			IssRectype:      "CDRNOTIF",
			IssRecaction:    "ADD",
			IssRecnum:       recNums.NextVal(),
			IssCompanyRegnr: "ARVD",
			KlLkeyAlias:     "",
			KlLKeyClr:       v.PAN,
			DbCdNotifSvcTyp: "SMSTXN",
			DbCdNotifTarget: v.Address,
		}
		jsonRec, err := json.Marshal(record)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b ADDCARD req marshaling record err: %v", err)
			return nil, err
		}
		recDetails.MdiRecords = append(recDetails.MdiRecords, jsonRec)
	}
	mdiDataJSON, err := json.Marshal(recDetails)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADD CRDNOTIF req marshaling err: %v", err)
		return nil, err
	}
	logger.Debugf("json ADD CRDNOTIF: %v", string(mdiDataJSON))

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/mdi", mdiDataJSON, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADD CRDNOTIF request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b ADD CRDNOTIF resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADD CRDNOTIF resp data unmarshal err: %v", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Data, &mdiData)
	if mdiData.Header.IRejMsg != "Approved" {
		return nil, fmt.Errorf("%s", mdiData.Header.IRejMsg)
	}
	if mdiData.Details != nil {
		if len(mdiData.Details) != 0 {
			return nil, fmt.Errorf("C_ACTIONCODE %s C_RSPCODE %s: Msg: %s", mdiData.Details[0].C_ACTIONCODE, mdiData.Details[0].C_RSPCODE, mdiData.Details[0].I_REJMSG)
		}
	}
	result = &resp
	return result, nil
}

func DeleteCardNotificationG2b(input models.MDIface) (result *d8corp.CommonResp, err error) {
	var (
		recDetails d8corp.MdiFile
		resp       d8corp.CommonResp
		mdiData    d8corp.MdiData
	)
	recNums := utils.NewSequence()

	for _, v := range input.GetRecords() {
		record := d8corp.MdiRecordDetails{
			IssRectype:      "CDRNOTIF",
			IssRecaction:    "DELETE",
			IssRecnum:       recNums.NextVal(),
			IssCompanyRegnr: "ARVD",
			KlLkeyAlias:     "",
			KlLKeyClr:       v.PAN,
			DbCdNotifSvcTyp: "SMSTXN",
			DbCdNotifTarget: v.Address,
		}
		jsonRec, err := json.Marshal(record)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b DELETECARDNOTIF req marshaling record err: %v", err)
			return nil, err
		}
		recDetails.MdiRecords = append(recDetails.MdiRecords, jsonRec)
	}
	mdiDataJSON, err := json.Marshal(recDetails)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b DELETECARDNOTIF req marshaling err: %v", err)
		return nil, err
	}
	logger.Debugf("json ADD CRDNOTIF: %v", string(mdiDataJSON))

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/mdi", mdiDataJSON, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b DELETECARDNOTIF request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b DELETECARDNOTIF resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b DELETECARDNOTIF resp data unmarshal err: %v", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Data, &mdiData)
	if mdiData.Header.IRejMsg != "Approved" {
		return nil, fmt.Errorf("%s", mdiData.Header.IRejMsg)
	}
	if mdiData.Details != nil {
		if len(mdiData.Details) != 0 {
			return nil, fmt.Errorf("C_ACTIONCODE %s C_RSPCODE %s: Msg: %s", mdiData.Details[0].C_ACTIONCODE, mdiData.Details[0].C_RSPCODE, mdiData.Details[0].I_REJMSG)
		}
	}
	result = &resp
	return result, nil
}
