package service

import (
	"converterapi/internal/config"
	"converterapi/internal/models"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func AddCardG2b(input models.MDIface) (resp interface{}, err error) {
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
		IssCompanyRegnr: "ARV",
		IssTimestamp:    "20230906120000123",
	}

	for i, v := range input.GetRecords() {
		separator := make([]byte, 0)
		if i != 0 {
			separator = json.RawMessage(",")
		}

		var firstSecret, firstName, lastName string
		if len(v.SecretInfo.Items) != 0 {
			firstSecret = v.SecretInfo.Items[0].Value
		}
		names := strings.Split(v.LatFIO, " ")
		if len(names) > 1 {
			lastName = names[0]
			firstName = names[1]
		}
		prior, _ := strconv.Atoi(v.MakePrior)
		record := d8corp.MdiRecordDetails{
			IssRectype:           "CARD",
			IssRecaction:         "ADD",
			IssRecnum:            recNums.NextVal(),
			IssCompanyRegnr:      "ARVD",
			IssCompanyRegnrAcc:   "ARVD",
			IssImpPvki:           1,
			DbCustomerCustcode:   firstSecret,
			DbCdproductCdproduct: "ARVDBT",
			DbAccountAccnum:      v.Account,
			DbAccountCurrcode:    v.CurrencyNo,
			KlLkeyAlias:          "93919",
			DbCardaCommCat:       "COM03",
			DbCardaEnroll3ds:     "1",
			DbCardaLimitCat:      "LIM01",
			DbCardEmbossname:     v.LatFIO,
			DbCardFirstname:      firstName,
			DbCardLastname:       lastName,
			DbCardMaidenname:     firstName,
			DbCrdaccPriority:     prior,
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

func AddPreissiedCardG2b(input models.MDIface) (resp interface{}, err error) {
	var (
		recDetails d8corp.MdiFile
	)
	recNums := utils.NewSequence()

	// filename := fmt.Sprintf("G2BISS-%v.JSON", time.Now().Local().Format("20060102-150405"))
	// header := d8corp.HeaderRecord{
	// 	IssRectype:      "HEADER",
	// 	IssRecaction:    "IMPORT",
	// 	CFilename:       filename, //"G2BISS-20060102-150405.JSON"
	// 	IssSourcesys:    "LK",
	// 	IssCompanyRegnr: "ARV",
	// 	IssTimestamp:    "20230906120000123",
	// }

	for _, v := range input.GetRecords() {
		// separator := make([]byte, 0)
		// if i != 0 {
		// 	separator = json.RawMessage(",")
		// }

		var firstSecret, firstName, lastName string
		if len(v.SecretInfo.Items) != 0 {
			firstSecret = v.SecretInfo.Items[0].Value
		}
		names := strings.Split(v.LatFIO, " ")
		if len(names) > 1 {
			lastName = names[0]
			firstName = names[1]
		}
		prior, _ := strconv.Atoi(v.MakePrior)
		record := d8corp.MdiRecordDetails{
			IssRectype:           "CARD",
			IssRecaction:         "ADD",
			IssRecnum:            recNums.NextVal(),
			IssCompanyRegnr:      "ARV",
			IssCompanyRegnrAcc:   "ARV",
			IssImpPvki:           1,
			DbCustomerCustcode:   firstSecret,
			DbCdproductCdproduct: "ARVDBT",
			//DbAccountAccnum:      v.Account,
			DbAccountCurrcode: v.CurrencyNo,
			KlLkeyAlias:       "505827",
			// DbCardaCommCat:    "COM03",
			DbCardaEnroll3ds: "1",
			DbCardaLimitCat:  "LIM01",
			DbCardEmbossname: v.LatFIO,
			DbCardFirstname:  firstName,
			DbCardLastname:   lastName,
			DbCardMaidenname: firstName,
			DbCrdaccPriority: prior,
		}
		jsonRec, err := json.Marshal(record)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b ADDCARD (preissued) req marshaling record err: %v", err)
			return nil, err
		}
		// recDetails.MdiRecords = append(recDetails.MdiRecords, separator)
		recDetails.MdiRecords = append(recDetails.MdiRecords, jsonRec)
	}

	// footer := d8corp.FooterRecord{
	// 	IssRectype:   "FOOTER",
	// 	IssRecaction: "IMPORT",
	// 	CFilename:    filename,
	// 	IssReccnt:    input.GetRecordsCount(),
	// }

	// headerJSON, err := json.Marshal(header)
	// if err != nil {
	// 	logger.Errorf("[SERVICE] D8 G2b ADDCARD (preissued) req marshaling err: %v", err)
	// 	return nil, err
	// }

	cardJSON, err := json.Marshal(recDetails)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD (preissued) req marshaling err: %v", err)
		return nil, err
	}
	logger.Infof("json ADDCARD: %v", string(cardJSON))

	// footerJSON, err := json.Marshal(footer)
	// if err != nil {
	// 	logger.Errorf("[SERVICE] D8 G2b ADDCARD (preissued) req marshaling err: %v", err)
	// 	return nil, err
	// }

	// mdiFile := d8corp.MdiFile{
	// 	MdiRecords: []json.RawMessage{
	// 		headerJSON,
	// 		cardJSON,
	// 		footerJSON,
	// 	},
	// }
	mdiDataJSON, err := json.Marshal(cardJSON)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD (preissued) req marshaling err: %v", err)
		return nil, err
	}

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/mdi", mdiDataJSON, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD (preissued) request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b ADDCARD (preissued) resp status: %v, body: %v", status, string(data))
	return data, nil
}
