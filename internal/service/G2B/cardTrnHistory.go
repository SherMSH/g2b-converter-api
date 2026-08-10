package service

import (
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
)

// historyPageSize - размер страницы при обходе истории операций.
// Процессинг ограничивает максимальный размер страницы своей настройкой и молча
// урезает ответ, поэтому запрашиваем заведомо безопасными порциями.
const historyPageSize = 100

// historyPageLimit - предохранитель от бесконечного обхода, если процессинг
// продолжит отдавать nextPagePresent при пустых страницах.
const historyPageLimit = 100

// GetCardTransactionHistory возвращает историю операций по карте.
//
// По спецификации (8.3) пагинация обязательна: в следующий запрос передаётся
// tlId последней полученной записи. Метод обходит страницы сам, пока не наберёт
// size записей или пока у процессинга не кончатся данные.
func GetCardTransactionHistory(pan, from, to string, size int) (transactions *d8corp.CardInfoData, err error) {
	if size <= 0 {
		size = historyPageSize
	}

	transactions = &d8corp.CardInfoData{}
	lastRetrievedId := -1

	for page := 0; page < historyPageLimit; page++ {
		want := size - len(transactions.CardTransactions)
		if want <= 0 {
			break
		}
		if want > historyPageSize {
			want = historyPageSize
		}

		batch, nextPage, err := getCardTransactionPage(pan, from, to, want, lastRetrievedId)
		if err != nil {
			// Первая страница - ошибка наружу; на последующих отдаём то, что успели
			// собрать: неполная выписка полезнее отказа.
			if page == 0 {
				return nil, err
			}
			logger.Warnf("[SERVICE] card history page %d err: %v", page, err)
			break
		}
		if len(batch) == 0 {
			break
		}

		transactions.CardTransactions = append(transactions.CardTransactions, batch...)
		lastRetrievedId = batch[len(batch)-1].TlId

		if !nextPage {
			break
		}
	}

	if len(transactions.CardTransactions) > size {
		transactions.CardTransactions = transactions.CardTransactions[:size]
	}
	return transactions, nil
}

// getCardTransactionPage запрашивает одну страницу истории.
// Второе возвращаемое значение - признак наличия следующей страницы.
func getCardTransactionPage(pan, from, to string, size, lastRetrievedId int) ([]d8corp.CardTransaction, bool, error) {
	var resp *d8corp.CommonResp

	req := d8corp.GetCardTrnHistoryReq{
		CardKey: d8corp.CardKey{
			Pan: pan,
		},
		DateLocalFrom: from,
		DateLocalTo:   to,
		PagingParams: d8corp.PagingParams{
			Size:            size,
			LastRetrievedId: lastRetrievedId,
		},
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b getCardTransactionHistory REQ marshaling err: %v", err)
		return nil, false, fmt.Errorf("[SERVICE] D8 G2b getCardTransactionHistory REQ marshaling err")
	}

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/kernel/1.0/getCardTransactionHistory", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b getCardTransactionHistory request sending err: %v", err)
		return nil, false, err
	}
	logger.Infof("[SERVICE] D8 G2b getCardTransactionHistory resp status: %v, body: %v (req %v)", status, string(data), string(jsonReq))

	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Errorf("[SERVICE] D8 G2b getCardTransactionHistory RESP marshaling err: %v", err)
		return nil, false, err
	}
	if resp.Status.Code != "0" {
		logger.Errorf("[SERVICE] D8 G2b getCardTransactionHistory RESP status %s", resp.Status.Code)
		return nil, false, fmt.Errorf("%s - %s", resp.Status.RspCode, resp.Status.Message)
	}

	page := &d8corp.CardInfoData{}
	if err = json.Unmarshal(resp.Data, page); err != nil {
		logger.Errorf("[SERVICE] D8 G2b getCardTransactionHistory DATA marshaling err: %v", err)
		return nil, false, err
	}
	return page.CardTransactions, resp.Paging.NextPagePresent, nil
}
