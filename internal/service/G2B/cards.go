package service

import (
	"converterapi/internal/utils"
	d8procweb "converterapi/pkg/d8-proc-web"
	"encoding/json"
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

	resp, err := d8procweb.D8procwebRequest(path, filters)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &foundCards)
	if err != nil {
		return nil, err
	}

	for i, v := range foundCards {
		cardInfo, _ := GetCardBasicInfo(v.LkeyID, "", utils.GetExpFormat4(v.Expdate))
		foundCards[i].PAN = cardInfo.CardBasicInfo.Lkey.Pan
		foundCards[i].StatCode = cardInfo.CardBasicInfo.StatCode
		foundCards[i].ProductType = cardInfo.CardBasicInfo.ProductType
	}

	return
}
