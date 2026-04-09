package service

import (
	"converterapi/internal/config"
	"converterapi/internal/models"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
	"time"
)

func AddCardNotificationG2b(input models.MDIface) (resp interface{}, err error) {
	var (
		recDetails d8corp.MdiFile
	)
	recNums := utils.NewSequence()

	filename := fmt.Sprintf("G2BISS-%v.JSON", time.Now().Local().Format("20060102-150405"))
	header := d8corp.HeaderRecord{
		IssRectype:      "HEADER",
		IssRecaction:    "IMPORT",
		CFilename:       filename, //"G2BISS-20060102-150405.JSON"
		IssSourcesys:    "LK",
		IssCompanyRegnr: "ARVD",
		IssTimestamp:    "20230906120000123",
	}

	for i, v := range input.GetRecords() {
		separator := make([]byte, 0)
		if i != 0 {
			separator = json.RawMessage(",")
		}
		record := d8corp.MdiRecordDetails{
			IssRectype:      "CDRNOTIF",
			IssRecaction:    "ADD",
			IssRecnum:       recNums.NextVal(),
			IssCompanyRegnr: "ARVD",
			KlLkeyAlias:     "",
			DbCdNotifSvcTyp: "SMSGEN",
			DbCdNotifTarget: v.Address,
		}
		jsonRec, err := json.Marshal(record)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b ADDCARD req marshaling record err: %v", err)
			return nil, err
		}
		recDetails.MdiRecords = append(recDetails.MdiRecords, separator)
		recDetails.MdiRecords = append(recDetails.MdiRecords, jsonRec)
	}

	footer := d8corp.FooterRecord{
		IssRectype:   "FOOTER",
		IssRecaction: "IMPORT",
		CFilename:    filename,
		IssReccnt:    input.GetRecordsCount(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD req marshaling err: %v", err)
		return nil, err
	}

	cardJSON, err := json.Marshal(recDetails)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD req marshaling err: %v", err)
		return nil, err
	}
	logger.Infof("json ADDCARD: %v", string(cardJSON))

	footerJSON, err := json.Marshal(footer)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD req marshaling err: %v", err)
		return nil, err
	}

	mdiFile := d8corp.MdiFile{
		MdiRecords: []json.RawMessage{
			headerJSON,
			cardJSON,
			footerJSON,
		},
	}
	mdiDataJSON, err := json.MarshalIndent(mdiFile, "", "  ")
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD req marshaling err: %v", err)
		return nil, err
	}

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/mdi", mdiDataJSON, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b ADDCARD resp status: %v, body: %v", status, string(data))
	return data, nil
}
