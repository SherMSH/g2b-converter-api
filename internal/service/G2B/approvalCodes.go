package service

import (
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/pkg/logger"
	"strconv"
	"sync"
)

// detailsWorkers ограничивает число одновременных запросов деталей транзакций.
// Процессинг обслуживает и онлайн-авторизации, поэтому выписка не должна
// занимать больше нескольких соединений разом.
const detailsWorkers = 5

// TxRef - ссылка на транзакцию для дозапроса деталей.
type TxRef struct {
	TlId      int
	EcTxRefno string
}

// GetTransactionDetailsBatch параллельно дозапрашивает детали транзакций.
//
// Нужен там, где список операций приходит без части полей: ни
// getCardTransactionHistory, ни getAccountStatement не возвращают код
// авторизации и реквизиты терминала, они есть только в getTransactionDetails.
//
// Ключ результата - tlId. Транзакции, по которым запрос не удался, в карту не
// попадают: пустое поле честнее подставленного неверного.
func GetTransactionDetailsBatch(refs []TxRef) map[int]d8corp.TransactionDetails {
	details := make(map[int]d8corp.TransactionDetails, len(refs))
	if len(refs) == 0 {
		return details
	}

	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)
	sem := make(chan struct{}, detailsWorkers)

	for _, ref := range refs {
		if ref.TlId == 0 {
			continue
		}
		wg.Add(1)
		go func(ref TxRef) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			trn, err := GetTransactionDetailsG2b(strconv.Itoa(ref.TlId), ref.EcTxRefno)
			if err != nil {
				logger.Warnf("[SERVICE] transaction details for tlId %d err: %v", ref.TlId, err)
				return
			}
			if trn == nil {
				return
			}

			mu.Lock()
			details[ref.TlId] = trn.Details
			mu.Unlock()
		}(ref)
	}
	wg.Wait()

	if len(details) != len(refs) {
		logger.Warnf("[SERVICE] details received for %d of %d transactions", len(details), len(refs))
	}
	return details
}

// GetApprovalCodes дозапрашивает коды авторизации по списку транзакций карты.
func GetApprovalCodes(trns []d8corp.CardTransaction) map[int]string {
	refs := make([]TxRef, 0, len(trns))
	for _, trn := range trns {
		refs = append(refs, TxRef{TlId: trn.TlId, EcTxRefno: trn.EcTxRefno})
	}

	details := GetTransactionDetailsBatch(refs)
	codes := make(map[int]string, len(details))
	for tlId, d := range details {
		if d.Aprvlcode != "" {
			codes[tlId] = d.Aprvlcode
		}
	}
	return codes
}
