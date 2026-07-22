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

func AddCardG2b(input models.MDIface, custExist, accExist bool) (mdiData *d8corp.MdiData, err error) {
	var (
		recDetails d8corp.MdiFile
		resp       d8corp.CommonResp
	)
	recNums := utils.NewSequence()
	for _, v := range input.GetRecords() {
		var secret, mobTel, firstName, lastName, nameInLat, lastNameinLat string
		if len(v.SecretInfo.Items) != 0 {
			secret = v.SecretInfo.Items[0].Value
		}
		names := strings.Split(strings.TrimSpace(v.LatFIO), " ")
		if len(names) > 1 {
			lastNameinLat = names[0]
			nameInLat = names[1]
		}
		fio := strings.Split(strings.TrimSpace(v.FIO), " ")
		if len(fio) > 2 {
			mobTel = fio[0]
			firstName = fio[2]
			lastName = fio[1]
		}
		bday, _ := time.Parse("02012006", v.Birthday)
		prior, _ := strconv.Atoi(v.MakePrior)
		prior++
		countryCode, er := strconv.Atoi(v.CountryRes)
		if er != nil {
			countryCode = 762
		}
		customerRec := d8corp.MdiRecordDetails{
			IssRectype:               "CUSTOMER",
			IssRecaction:             "ADD",
			IssRecnum:                recNums.NextVal(),
			IssCompanyRegnr:          "ARV",
			DbCustomerTypeCode:       "0",
			DbCustomerCustcode:       v.ExtID,
			DbCustomerFirstName:      firstName,
			DbCustomerLastName:       lastName,
			DbCustomerLatinFirstName: nameInLat,
			DbCustomerLatinLastName:  lastNameinLat,
			DbCustomerDateBirth:      bday.Format("20060102"),
			DbCustomerHomeCountry:    countryCode,
			DbCustomerPassPhrase:     secret,
			DbCustomerMobTel:         mobTel,
			DbCustomerDocument:       v.PasNom,
		}
		accRec := d8corp.MdiRecordDetails{
			IssRectype:         "ACCOUNT",
			IssRecaction:       "ADD",
			IssRecnum:          recNums.NextVal(),
			IssCompanyRegnr:    "ARV",
			DbCustomerCustcode: v.ExtID,
			DbAccountCurrcode:  v.CurrencyNo,
			DbAccountAccnum:    v.Account,
			DbAccountTypecode:  "00",
		}
		cardRec := d8corp.MdiRecordDetails{
			IssRectype:      "CARD",
			IssRecaction:    "ADD",
			IssRecnum:       recNums.NextVal(),
			IssCompanyRegnr: "ARV",
			// IssCompanyRegnrAcc:   "ARVD",
			IssImpPvki:           1,
			DbCustomerCustcode:   v.ExtID,
			DbCdproductCdproduct: "ARVDBT",
			DbAccountAccnum:      v.Account,
			DbAccountCurrcode:    v.CurrencyNo,
			DbCardaCommCat:       "COM03",
			DbCardaEnroll3ds:     "1",
			DbCardaLimitCat:      "LIM01",
			DbCardEmbossname:     v.NameOnCard,
			DbCardFirstname:      nameInLat,
			DbCardLastname:       lastNameinLat,
			DbCardMaidenname:     lastNameinLat,
			DbCrdaccPriority:     prior,
		}
		// linkRec := d8corp.MdiRecordDetails{
		// 	IssRectype: "CRDACC",
		// 	IssRecaction: "ADD",
		// 	IssRecnum: recNums.NextVal(),
		// 	IssCompanyRegnr: "ARV",
		// 	KlLKeyClr: "",
		// 	KlLkeySeqno: 0,
		// 	DbAccountAccnum: v.Account,
		// 	DbAccountCurrcode: v.CurrencyNo,
		// 	DbCrdaccPriority: prior,
		// }

		jsonCustmr, err := json.Marshal(customerRec)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b ADDCARD customer req marshaling record err: %v", err)
			return nil, err
		}
		jsonAccnt, err := json.Marshal(accRec)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b ADDCARD account req marshaling record err: %v", err)
			return nil, err
		}
		jsonCrd, err := json.Marshal(cardRec)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b ADDCARD card req marshaling record err: %v", err)
			return nil, err
		}
		if !custExist {
			recDetails.MdiRecords = append(recDetails.MdiRecords, jsonCustmr)
		}
		if !accExist {
			recDetails.MdiRecords = append(recDetails.MdiRecords, jsonAccnt)
		}
		recDetails.MdiRecords = append(recDetails.MdiRecords, jsonCrd)
	}

	// footer := d8corp.FooterRecord{
	// 	IssRectype:   "FOOTER",
	// 	IssRecaction: "IMPORT",
	// 	CFilename:    filename,
	// 	IssReccnt:    input.GetRecordsCount(),
	// }

	// headerJSON, err := json.Marshal(header)
	// if err != nil {
	// 	logger.Errorf("[SERVICE] D8 G2b ADDCARD req marshaling err: %v", err)
	// 	return nil, err
	// }

	cardJSON, err := json.Marshal(recDetails)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD req marshaling err: %v", err)
		return nil, err
	}
	logger.Infof("json ADDCARD: %v", string(cardJSON))

	// footerJSON, err := json.Marshal(footer)
	// if err != nil {
	// 	logger.Errorf("[SERVICE] D8 G2b ADDCARD req marshaling err: %v", err)
	// 	return nil, err
	// }

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/mdi", cardJSON, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b ADDCARD resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD common RESP unmarshaling err: %v", err)
		return nil, err
	}
	if resp.Status.Code != "0" {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD RESP status %s", resp.Status.Code)
		return nil, fmt.Errorf("%s - %s", resp.Status.RspCode, resp.Status.Message)
	}

	mdiData = new(d8corp.MdiData)
	err = json.Unmarshal(resp.Data, mdiData)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD mdiData marshaling err: %v", err)
		return nil, err
	}

	return mdiData, nil
}

func AddPreissiedCardG2b(input models.MDIface) (mdiData *d8corp.MdiData, err error) {
	var (
		recDetails d8corp.MdiFile
		resp       d8corp.CommonResp
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
		prior++
		record := d8corp.MdiRecordDetails{
			IssRectype:           "CARD",
			IssRecaction:         "ADD",
			IssRecnum:            recNums.NextVal(),
			IssCompanyRegnr:      "ARV",
			IssCompanyRegnrAcc:   "ARV",
			IssImpPvki:           1,
			DbCustomerCustcode:   firstSecret,
			DbCdproductCdproduct: "ARVDBT",
			DbAccountAccnum:      v.Account,
			DbAccountCurrcode:    v.CurrencyNo,
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

	mdiDataJSON, err := json.Marshal(recDetails)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD (preissued) req marshaling err: %v", err)
		return nil, err
	}
	logger.Infof("json ADDCARD (preissued): %v", string(mdiDataJSON))

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
	// mdiDataJSON, err := json.Marshal(cardJSON)
	// if err != nil {
	// 	logger.Errorf("[SERVICE] D8 G2b ADDCARD (preissued) req marshaling err: %v", err)
	// 	return nil, err
	// }

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/mdi", mdiDataJSON, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD (preissued) request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b ADDCARD (preissued) resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD RESP marshaling err: %v", err)
		return nil, err
	}
	if resp.Status.Code != "0" {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD RESP status %s", resp.Status.Code)
		return nil, fmt.Errorf("%s - %s", resp.Status.RspCode, resp.Status.Message)
	}

	err = json.Unmarshal(resp.Data, mdiData)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD RESP marshaling err: %v", err)
		return nil, err
	}

	return mdiData, nil
}
