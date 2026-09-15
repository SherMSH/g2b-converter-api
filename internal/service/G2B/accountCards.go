package service

import (
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
)

// accountCardsPageSize - размер страницы при выборке карт счёта.
// У одного счёта карт немного, но пагинация в запросе обязательна.
const accountCardsPageSize = 100

// GetAccountCardListG2b возвращает карты, привязанные к счёту
// (xmiss/getAccountCardList, 7.28).
//
// Это не то же самое, что getCardList по клиенту: у клиента может быть
// несколько счетов, и карты других его счетов к этому отношения не имеют.
// Ответ сразу содержит полные данные карты, поэтому дозапрашивать
// getCardBasicInfo по каждой карте не нужно.
func GetAccountCardListG2b(accNum, currency string) (cards []d8corp.CardBasicInfo, err error) {
	var resp *d8corp.CommonResp

	req := d8corp.GetAccountCardListReq{
		AccountKey: d8corp.AccountKey{
			AccountNumber: accNum,
			Currency:      currency,
		},
		PagingAdvanced: d8corp.PagingAdvanced{
			Size: accountCardsPageSize,
			Page: 1,
		},
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b getAccountCardList REQ marshaling err: %v", err)
		return nil, fmt.Errorf("[SERVICE] D8 G2b getAccountCardList REQ marshaling err")
	}

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/getAccountCardList", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b getAccountCardList request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b getAccountCardList resp status: %v, body: %v (req %v)", status, string(data), string(jsonReq))

	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Errorf("[SERVICE] D8 G2b getAccountCardList RESP marshaling err: %v", err)
		return nil, err
	}
	if resp.Status.Code != "0" {
		logger.Errorf("[SERVICE] D8 G2b getAccountCardList RESP status %+v", resp.Status)
		return nil, StatusError(resp.Status)
	}

	list := d8corp.CardList{}
	if err = json.Unmarshal(resp.Data, &list); err != nil {
		logger.Errorf("[SERVICE] D8 G2b getAccountCardList DATA marshaling err: %v", err)
		return nil, err
	}
	return list.Cards, nil
}
