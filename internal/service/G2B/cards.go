package service

import (
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	d8procweb "converterapi/pkg/d8-proc-web"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
)

func GetCardsListG2b(custcode, currcode string) (foundCards []d8procweb.CardData, err error) {
	path := "/api/miss/v1/getCardList"
	filters := []d8procweb.Filter{
		{
			Column: "custcode",
			Values: []string{custcode},
		},
		{
			Column: "currcode",
			Values: []string{currcode},
		},
	}

	resp, err := d8procweb.Request(path, filters)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &foundCards)
	if err != nil {
		return nil, err
	}

	for i, v := range foundCards {
		// В CardKey допустим ровно один набор реквизитов (5.3.5): либо lkeyId,
		// либо PAN со сроком действия. Мы передавали lkeyId вместе со сроком, и
		// при расхождении срока процессинг отвечал "Card not found" - карта
		// оставалась без PAN и статуса.
		cardInfo, err := GetCardBasicInfo(v.LkeyID, "", "")
		if err != nil {
			logger.Warnf("[SERVICE] GetCardsList: данные карты lkeyId %d не получены: %v", v.LkeyID, err)
			continue
		}
		if cardInfo != nil {
			foundCards[i].PAN = cardInfo.CardBasicInfo.Lkey.Pan
			foundCards[i].StatCode = cardInfo.CardBasicInfo.StatCode
			foundCards[i].ProductType = cardInfo.CardBasicInfo.ProductType
		}
	}

	return
}

func GetExpDateByPan(pan string) (expdate string, err error) {
	var (
		resp     *d8corp.CommonResp
		cardList *d8corp.CardList
	)
	req := d8corp.GetCardInfoReq{
		CardKey: d8corp.CardKey{
			Pan: pan,
		},
	}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b GetCVVG2b REQ marshaling err: %v", err)
		return "", fmt.Errorf("[SERVICE] D8 G2b CardList REQ marshaling err")
	}
	logger.Infof("[SERVICE] D8 G2b CardList REQ %v", string(jsonReq))
	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/getCardListByLkey", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b CardList request sending err: %v", err)
		return "", err
	}
	logger.Infof("[SERVICE] D8 G2b CardList resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b CardList RESP marshaling err: %v", err)
		return "", err
	}
	if resp.Status.Code != "0" {
		logger.Errorf("[SERVICE] D8 G2b CardList RESP status %s", resp.Status.Code)
		return "", fmt.Errorf("%s - %s", resp.Status.RspCode, resp.Status.Message)
	}

	err = json.Unmarshal(resp.Data, &cardList)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b CardList DATA marshaling err: %v", err)
		return "", err
	}

	expdates := make([]string, 0)
	for _, card := range cardList.Cards {
		expdates = append(expdates, card.ExpiryDate)
	}
	if len(expdates) == 0 {
		return "", fmt.Errorf("card list is empty!")
	}
	expdate = utils.ConvertYYYYMMDDtoYYMM(expdates[0])
	return expdate, nil
}
