package createcustomerandaccount

import (
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"strings"
	"time"
)

func CreateCustomersAndAccountsG2b(input Root) (resp interface{}, err error) {
	var (
		recDetails d8corp.MdiFile
	)
	recNums := utils.NewSequence()
	for _, v := range input.Records {
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
		recDetails.MdiRecords = append(recDetails.MdiRecords, jsonRec)

		jsonRec, err = json.Marshal(accountRec)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b add account req marshaling record err: %v", err)
			return nil, err
		}
		recDetails.MdiRecords = append(recDetails.MdiRecords, jsonRec)
	}
	reqJSON, err := json.Marshal(recDetails)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b CreateCustomersAndAccounts req marshaling err: %v", err)
		return nil, err
	}
	logger.Infof("json ADD CUSTOMER and ACCOUNT: %v", string(reqJSON))

	mdiFile := d8corp.MdiFile{
		MdiRecords: []json.RawMessage{
			reqJSON,
		},
	}
	mdiDataJSON, err := json.MarshalIndent(mdiFile, "", "  ")
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b CreateCustomersAndAccounts req marshaling err: %v", err)
		return nil, err
	}

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/mdi", mdiDataJSON, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADD CUSTOMER and ACCOUNT request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b ADD CUSTOMER and ACCOUNT resp status: %v, body: %v", status, string(data))
	return data, nil
}

// func CreateCustomersG2b() error {
// 	return nil
// }
