package posrequestrq

import (
	"fmt"
	"sync"
	"time"
)

// dedupTTL - окно, в котором повторный запрос считается дублем.
//
// Повтор при обрыве связи или таймауте приходит в пределах секунд, поэтому
// нескольких минут достаточно. Больше брать опасно: если партнёр переиспользует
// номера операций, мы начнём глушить настоящие платежи.
const dedupTTL = 5 * time.Minute

// posDedup хранит ответы по номеру операции партнёра.
//
// Собственной БД у сервиса нет, поэтому память процесса: перезапуск окно
// сбрасывает, а при нескольких экземплярах каждый считает дубли сам. Для защиты
// от повторов внутри одной сессии связи этого достаточно, для настоящей
// идемпотентности нужно общее хранилище.
type posDedup struct {
	mu      sync.Mutex
	entries map[string]dedupEntry
}

type dedupEntry struct {
	resp Envelope
	at   time.Time
}

var dedupCache = posDedup{entries: make(map[string]dedupEntry)}

// dedupKey - ключ дубля.
//
// Одного номера операции мало: партнёры нумеруют операции по-своему и номера
// повторяются между днями и терминалами. Поэтому в ключ входят ещё карта, код
// операции и сумма - совпадение всего набора означает буквально тот же платёж.
func dedupKey(req Request) string {
	if req.TranNumber == "" {
		return ""
	}
	return fmt.Sprintf("%s|%s|%d|%.2f|%s", req.TranNumber, req.PAN, req.TranCode, req.Amount, req.PAN2)
}

// get возвращает ответ на такой же запрос, если он был недавно.
func (d *posDedup) get(key string) (Envelope, bool) {
	if key == "" {
		return Envelope{}, false
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	entry, ok := d.entries[key]
	if !ok {
		return Envelope{}, false
	}
	if time.Since(entry.at) > dedupTTL {
		delete(d.entries, key)
		return Envelope{}, false
	}
	return entry.resp, true
}

// put запоминает ответ. Заодно выбрасывает протухшие записи - отдельная уборка
// не нужна, запросов достаточно, чтобы карта не росла.
func (d *posDedup) put(key string, resp *Envelope) {
	if key == "" || resp == nil {
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	for k, entry := range d.entries {
		if time.Since(entry.at) > dedupTTL {
			delete(d.entries, k)
		}
	}
	d.entries[key] = dedupEntry{resp: *resp, at: time.Now()}
}
