package service

import (
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/pkg/logger"
	"strconv"
	"sync"
)

// approvalWorkers ограничивает число одновременных запросов деталей транзакций.
// Процессинг обслуживает и онлайн-авторизации, поэтому выписка не должна
// занимать больше нескольких соединений разом.
const approvalWorkers = 5

// GetApprovalCodes дозапрашивает коды авторизации по списку транзакций.
//
// getCardTransactionHistory не возвращает aprvlcode, он есть только в ответе
// getTransactionDetails - поэтому на каждую строку выписки нужен свой запрос.
// Запросы идут параллельно с ограничением approvalWorkers.
//
// Ключ результата - tlId транзакции. Строки, по которым details получить не
// удалось, в карту не попадают: пустой код честнее подставленного неверного.
func GetApprovalCodes(trns []d8corp.CardTransaction) map[int]string {
	codes := make(map[int]string, len(trns))
	if len(trns) == 0 {
		return codes
	}

	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)
	sem := make(chan struct{}, approvalWorkers)

	for _, trn := range trns {
		wg.Add(1)
		go func(tlId int, ecTxRefno string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			details, err := GetTransactionDetailsG2b(strconv.Itoa(tlId), ecTxRefno)
			if err != nil {
				logger.Warnf("[SERVICE] approval code for tlId %d err: %v", tlId, err)
				return
			}
			if details == nil || details.Details.Aprvlcode == "" {
				return
			}

			mu.Lock()
			codes[tlId] = details.Details.Aprvlcode
			mu.Unlock()
		}(trn.TlId, trn.EcTxRefno)
	}
	wg.Wait()

	if len(codes) != len(trns) {
		logger.Warnf("[SERVICE] approval codes received for %d of %d transactions", len(codes), len(trns))
	}
	return codes
}
