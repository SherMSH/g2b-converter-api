package service

import (
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
)

// GetAccountStatementG2b запрашивает выписку по счёту - xmiss/getAccountStatement (7.30).
//
// В отличие от getCardTransactionHistory отдаёт движения самого счёта, включая
// операции, прошедшие мимо карты, и состояние баланса на начало и конец периода.
//
// Даты в формате YYYYMMDD. Пагинация - size + page, нумерация страниц с 1.
//
// Запрашиваем финансовые движения вместе с блокировками. Без logType процессинг
// отдаёт весь журнал, включая правки карточки счёта: баланс по ним не менялся,
// и в выписке они превращаются в строки с нулевой суммой.
//
// Именно 3, а не 1: блокировка денег со счёта не снимает, но доступный остаток
// уменьшает - для держателя это движение, и в выписке оно нужно. С logType=1
// преавторизации и холды пропали бы молча.
func GetAccountStatementG2b(accNum, currency, dateFrom, dateTo string, size, page int) (statement *d8corp.AccountStatementData, err error) {
	var resp *d8corp.CommonResp

	if page < 1 {
		page = 1
	}
	req := d8corp.GetAccountStatementReq{
		AccountKey: d8corp.AccountKey{
			AccountNumber: accNum,
			Currency:      currency,
		},
		DateFrom: dateFrom,
		DateTo:   dateTo,
		LogType:  d8corp.LogTypeFinBlk,
		PagingAdvanced: d8corp.PagingAdvanced{
			Size: size,
			Page: page,
		},
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b getAccountStatement REQ marshaling err: %v", err)
		return nil, fmt.Errorf("[SERVICE] D8 G2b getAccountStatement REQ marshaling err")
	}

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/miss/1.0/getAccountStatement", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b getAccountStatement request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b getAccountStatement resp status: %v, body: %v (req %v)", status, string(data), string(jsonReq))

	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Errorf("[SERVICE] D8 G2b getAccountStatement RESP marshaling err: %v", err)
		return nil, err
	}
	if resp.Status.Code != "0" {
		logger.Errorf("[SERVICE] D8 G2b getAccountStatement RESP status %s", resp.Status.Code)
		return nil, fmt.Errorf("%s - %s", resp.Status.RspCode, resp.Status.Message)
	}
	if err = json.Unmarshal(resp.Data, &statement); err != nil {
		logger.Errorf("[SERVICE] D8 G2b getAccountStatement DATA marshaling err: %v", err)
		return nil, err
	}
	if statement == nil {
		return nil, fmt.Errorf("no data")
	}
	return statement, nil
}
