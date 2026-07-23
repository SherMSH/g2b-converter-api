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

func AddCardG2b(input models.MDIface) (mdiData *d8corp.MdiData, err error) {
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
		// Вызов функции исполнения MDI запроса
		err = CreateCustomerTry(&customerRec)
		if err != nil {
			err = nil
			continue
		}

		// Вызов функции исполнения MDI запроса
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
		err = CreateAccountTry(&accRec)
		if err != nil {
			err = nil
			continue
		}
		cardRec := d8corp.MdiRecordDetails{
			IssRectype:           "CARD",
			IssRecaction:         "ADD",
			IssRecnum:            recNums.NextVal(),
			IssCompanyRegnr:      "ARV",
			IssImpPvki:           3,
			IssGenPin:            1,
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
		jsonCrd, err := json.Marshal(cardRec)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b ADDCARD card req marshaling record err: %v", err)
			return nil, err
		}
		recDetails.MdiRecords = append(recDetails.MdiRecords, jsonCrd)
	}

	cardJSON, err := json.Marshal(recDetails)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD req marshaling err: %v", err)
		return nil, err
	}
	logger.Infof("json ADDCARD: %v", string(cardJSON))

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

	for _, v := range input.GetRecords() {
		var nameInLat, lastNameinLat string
		secret := "imtiyoz"
		if len(v.SecretInfo.Items) != 0 {
			secret = v.SecretInfo.Items[0].Value
		}
		names := strings.Split(strings.TrimSpace(v.LatFIO), " ")
		if len(names) > 1 {
			lastNameinLat = names[0]
			nameInLat = names[1]
		}

		prior, _ := strconv.Atoi(v.MakePrior)
		prior++
		record := d8corp.MdiRecordDetails{
			IssRectype:           "CARD",
			IssRecaction:         "ADD",
			IssRecnum:            recNums.NextVal(),
			IssCompanyRegnr:      "ARV",
			IssCompanyRegnrAcc:   "ARV",
			IssImpPvki:           3,
			IssGenPin:            1,
			DbCustomerCustcode:   secret,
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
		jsonRec, err := json.Marshal(record)
		if err != nil {
			logger.Errorf("[SERVICE] D8 G2b ADDCARD (preissued) req marshaling record err: %v", err)
			return nil, err
		}
		recDetails.MdiRecords = append(recDetails.MdiRecords, jsonRec)
	}

	mdiDataJSON, err := json.Marshal(recDetails)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD (preissued) req marshaling err: %v", err)
		return nil, err
	}
	logger.Infof("json ADDCARD (preissued): %v", string(mdiDataJSON))

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/mdi", mdiDataJSON, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD (preissued) request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b ADDCARD (preissued) resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD (preissued) RESP marshaling err: %v", err)
		return nil, err
	}
	if resp.Status.Code != "0" {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD (preissued) RESP status %s", resp.Status.Code)
		return nil, fmt.Errorf("%s - %s", resp.Status.RspCode, resp.Status.Message)
	}

	mdiData = new(d8corp.MdiData)
	err = json.Unmarshal(resp.Data, mdiData)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD (preissued) mdiData marshaling err: %v", err)
		return nil, err
	}

	return mdiData, nil
}
