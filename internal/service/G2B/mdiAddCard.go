package service

import (
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	createcardsout "converterapi/internal/models/OFFLINE/CreateCardsOut"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"strings"
)

func AddCardG2b(input createcardsout.Root) (resp interface{}, err error) {
	var (
		req4D8 d8corp.Request
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
		record := d8corp.MdiRecord{
			IssRectype:           "CARD",
			IssRecaction:         "ADD",
			IssRecnum:            recNums.NextVal(),
			IssCompanyRegnr:      "ARVD",
			IssCompanyRegnrAcc:   "ARVD",
			IssImpPvki:           "1",
			DbCustomerCustcode:   firstSecret,
			DbCdproductCdproduct: "CD01",
			DbAccountAccnum:      v.Account,
			DbAccountCurrcode:    v.CurrencyNo,
			KlLkeyAlias:          "93919",
			KlLkeySeqno:          "0",
			DbCardaExpdate:       v.CancelDate,
			DbCardaCommCat:       "COM03",
			DbCardaEnroll3ds:     "1",
			DbCardaLimitCat:      "LIM01",
			DbCardEmbossname:     v.LatFIO,
			DbCardFirstname:      firstName,
			DbCardLastname:       lastName,
			DbCardMaidenname:     firstName,
			DbCardDeliveryPoint:  "000",
			DbCrdaccPriority:     v.MakePrior,
		}
		req4D8.MDIRecords = append(req4D8.MDIRecords, record)
	}

	reqBody, err := json.Marshal(req4D8)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD req marshaling err: %v", err)
		return nil, err
	}
	logger.Infof("json ADDCARD: %v", string(reqBody))

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/mdi", reqBody, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b ADDCARD request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b ADDCARD resp status: %v, body: %v", status, string(data))
	return data, nil
}
