package createcustomerandaccount

import (
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func CreateCustomersAndAccountsG2b(input Root) (resp interface{}, err error) {
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
		IssCompanyRegnr: "COMPANY1",
		IssTimestamp:    "20230906120000123",
	}

	for i, v := range input.Records {
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
		bday, _ := time.Parse("02012006", v.Birthday)
		customerRec := d8corp.MdiRecordDetails{
			IssRectype:            "CUSTOMER",
			IssRecaction:          "ADD",
			IssRecnum:             recNums.NextVal(),
			IssCompanyRegnr:       "ARVD",
			DbCustomerTypeCode:    0,
			DbCustomerCustcode:    firstSecret,
			DbCustomerFirstName:   firstName,
			DbCustomerLastName:    lastName,
			DbCustomerDateBirth:   bday.Format("20060102"),
			DbCustomerHomeCountry: v.CountryRes,
		}

		accountRec := d8corp.MdiRecordDetails{
			IssRectype:         "ACCOUNT",
			IssRecaction:       "ADD",
			IssRecnum:          recNums.NextVal(),
			IssCompanyRegnr:    "ARVD",
			DbCustomerCustcode: firstSecret,
			DbAccountCurrcode:  v.Currency,
			DbAccountAccnum:    v.Account,
			DbAccountTypecode:  "00",
		}

		jsonRec, err := json.Marshal(customerRec)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b add customer req marshaling record err: %v", err)
			return nil, err
		}

		recDetails.MdiRecords = append(recDetails.MdiRecords, separator)
		jsonRec, err = json.Marshal(accountRec)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b add account req marshaling record err: %v", err)
			return nil, err
		}

		recDetails.MdiRecords = append(recDetails.MdiRecords, separator)
		recDetails.MdiRecords = append(recDetails.MdiRecords, jsonRec)
	}

	footer := d8corp.FooterRecord{
		IssRectype:   "FOOTER",
		IssRecaction: "IMPORT",
		CFilename:    filename,
		IssReccnt:    2 * len(input.Records),
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

// func CreateCustomersG2b() error {
// 	return nil
// }
