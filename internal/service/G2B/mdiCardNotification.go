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
		smsTxnRec := d8corp.MdiRecordDetails{
			IssRectype:      "CDRNOTIF",
			IssRecaction:    "ADD",
			IssRecnum:       recNums.NextVal(),
			IssCompanyRegnr: "ARVD",
			KlLkeyAlias:     "",
			KlLKeyClr:       v.PAN,
			DbCdNotifSvcTyp: "SMSTXN",
			DbCdNotifTarget: v.Address,
		}
		smsGenRec := d8corp.MdiRecordDetails{
			IssRectype:      "CDRNOTIF",
			IssRecaction:    "ADD",
			IssRecnum:       recNums.NextVal(),
			IssCompanyRegnr: "ARVD",
			KlLkeyAlias:     "",
			KlLKeyClr:       v.PAN,
			DbCdNotifSvcTyp: "SMSGEN",
			DbCdNotifTarget: v.Address,
		}

		jsonRec1, err := json.Marshal(smsTxnRec)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b ADDCARD SMSTXN req marshaling record err: %v", err)
			return nil, err
		}
		jsonRec2, err := json.Marshal(smsGenRec)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b ADDCARD SMSGEN req marshaling record err: %v", err)
			return nil, err
		}
		recDetails.MdiRecords = append(recDetails.MdiRecords, jsonRec1, jsonRec2)
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
		smsTxnRec := d8corp.MdiRecordDetails{
			IssRectype:      "CDRNOTIF",
			IssRecaction:    "DELETE",
			IssRecnum:       recNums.NextVal(),
			IssCompanyRegnr: "ARVD",
			KlLkeyAlias:     "",
			KlLKeyClr:       v.PAN,
			DbCardaExpdate:  v.ExpDate,
			DbCdNotifSvcTyp: "SMSTXN",
			DbCdNotifTarget: v.Address,
		}
		smsGenRec := d8corp.MdiRecordDetails{
			IssRectype:      "CDRNOTIF",
			IssRecaction:    "DELETE",
			IssRecnum:       recNums.NextVal(),
			IssCompanyRegnr: "ARVD",
			KlLkeyAlias:     "",
			KlLKeyClr:       v.PAN,
			DbCardaExpdate:  v.ExpDate,
			DbCdNotifSvcTyp: "SMSGEN",
			DbCdNotifTarget: v.Address,
		}

		jsonRec1, err := json.Marshal(smsTxnRec)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b DELETECARDNOTIF SMSTXN req marshaling record err: %v", err)
			return nil, err
		}
		jsonRec2, err := json.Marshal(smsGenRec)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b DELETECARDNOTIF SMSGEN req marshaling record err: %v", err)
			return nil, err
		}
		recDetails.MdiRecords = append(recDetails.MdiRecords, jsonRec1, jsonRec2)
	}

	mdiDataJSON, err := json.Marshal(recDetails)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b DELETECARDNOTIF req marshaling err: %v", err)
		return nil, err
	}
	logger.Debugf("json DELETE CRDNOTIF: %v", string(mdiDataJSON))

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
