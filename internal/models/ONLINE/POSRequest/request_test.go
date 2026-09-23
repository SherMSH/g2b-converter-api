package posrequestrq

import (
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"strings"
	"testing"
	"time"
)

func TestGetTxnType(t *testing.T) {
	tests := []struct {
		name string
		req  Request
		want utils.TxnType
	}{
		{
			name: "покупка",
			req:  Request{TranCode: Debit},
			want: utils.Sales,
		},
		{
			name: "зачисление",
			req:  Request{TranCode: Credit},
			want: utils.Deposit,
		},
		{
			name: "P2P с картой получателя - карта в карту",
			req:  Request{TranCode: P2P, PAN2: "5058270530003879"},
			want: utils.C2C,
		},
		{
			name: "перевод на счёт - карта в счёт",
			req:  Request{TranCode: Transfer, ToAccount: "20216972300001176308"},
			want: utils.C2A,
		},
		{
			name: "перевод со счёта на карту",
			req:  Request{TranCode: Transfer, PAN2: "5058270530003879", FromAccount: "20216972300001176308"},
			want: utils.A2C,
		},
		{
			name: "перевод без реквизитов получателя не распознаётся",
			req:  Request{TranCode: P2P},
			want: "",
		},
		{
			name: "неизвестный код не превращается в списание",
			req:  Request{TranCode: 999},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.req.GetTxnType(); got != tt.want {
				t.Errorf("GetTxnType() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetTxnTypeCheckCard(t *testing.T) {
	req := Request{TranCode: CheckCard, CVV2: "123"}
	if got := req.GetTxnType(); got != utils.Accver {
		t.Errorf("проверка карты должна давать ACCVER, а даёт %q", got)
	}
	if req.GetCvv2() != "123" {
		t.Errorf("CVV2 не отдаётся: %q", req.GetCvv2())
	}
}

// Процессинг часто отвечает общим отказом (Do not honour -> 50). По таблице
// партнёра код должен определяться статусом карты: скомпрометирована -> 75.
func TestRefineDeclineCode(t *testing.T) {
	card := func(statCode, acctStatus string) *d8corp.CardInfoData {
		return &d8corp.CardInfoData{
			CardBasicInfo: d8corp.CardBasicInfo{StatCode: statCode},
			CardAccounts:  []d8corp.CardAccount{{StatCode: acctStatus}},
		}
	}

	tests := []struct {
		name    string
		code    string
		card    *d8corp.CardInfoData
		isDebit bool
		want    string
	}{
		{"общий отказ по скомпрометированной карте", "50", card("08", "00"), true, "75"},
		{"мошенническое использование тоже 75", "50", card("16", "00"), true, "75"},
		{"неизвестный отказ по утерянной карте", "68", card("12", "00"), true, "40"},
		{"украденная карта", "50", card("13", "00"), false, "41"},
		{"истёкшая карта", "50", card("11", "00"), true, "51"},
		{"ограниченная карта, расход", "50", card("10", "00"), true, "58"},
		{"ограниченная карта, зачисление - статус не мешает", "50", card("10", "00"), false, "50"},
		{"закрытый счёт при активной карте", "50", card("00", "09"), true, "56"},
		{"только приход, расход", "50", card("00", "02"), true, "55"},
		{"статус важнее кода процессинга", "62", card("10", "00"), true, "58"},
		{"скомпрометированная карта важнее нехватки средств", "59", card("08", "00"), true, "75"},
		{"при исправной карте код процессинга сохраняется", "59", card("00", "00"), true, "59"},
		{"неверный PIN при исправной карте", "53", card("00", "00"), true, "53"},
		{"без данных карты код остаётся прежним", "50", nil, true, "50"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := refineDeclineCode(tt.code, tt.card, tt.isDebit); got != tt.want {
				t.Errorf("refineDeclineCode(%q) = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

// Связанные операции берутся из transactionGroups процессинга: тип 2 - реверс,
// тип 5 - части P2P-перевода. Сама операция в список не попадает.
func TestRelatedTran(t *testing.T) {
	details := &d8corp.Transaction{Details: d8corp.TransactionDetails{
		TlId: 1500,
		TransactionGroups: []d8corp.TransactionGroup{{
			GroupId:   142,
			GroupType: d8corp.GroupTypeReversalLinkage,
			Transactions: []d8corp.TransactionBasic{
				{TlId: 1500, TxnCode: 0, ActionCode: "0", RspCode: "00"},
				{TlId: 1504, TxnCode: 21, ActionCode: "0", RspCode: "00"},
			},
		}},
	}}

	got := relatedTran(details)
	if len(got.Rows) != 1 {
		t.Fatalf("ожидалась одна связанная операция, получено %d", len(got.Rows))
	}
	row := got.Rows[0]
	if row.RelatedTranId != "1504" {
		t.Errorf("RelatedTranId = %q, want 1504", row.RelatedTranId)
	}
	if row.RelatedTranCode != "140" {
		t.Errorf("RelatedTranCode = %q, want 140 (код операции переведён в кодировку TWO)", row.RelatedTranCode)
	}
	if row.RelatedAuthRespCode != "1" {
		t.Errorf("RelatedAuthRespCode = %q, want 1", row.RelatedAuthRespCode)
	}
}

func TestRelatedTranEmpty(t *testing.T) {
	if rows := relatedTran(nil).Rows; rows != nil {
		t.Errorf("без деталей операции связей быть не должно: %+v", rows)
	}
	details := &d8corp.Transaction{Details: d8corp.TransactionDetails{TlId: 1}}
	if rows := relatedTran(details).Rows; rows != nil {
		t.Errorf("без групп связей быть не должно: %+v", rows)
	}
}

// Плечи перевода процессинг отклоняет дежурным "AFT Rejected" - причину для
// партнёра в этом случае собираем сами.
func TestIsPlaceholderReason(t *testing.T) {
	for _, reason := range []string{"", "  ", "AFT Rejected", "aft rejected", "OCT Rejected"} {
		if !isPlaceholderReason(reason) {
			t.Errorf("%q - заглушка процессинга, причину нужно собрать самим", reason)
		}
	}
	real := "PIN tries exceeded has been set on card status"
	if isPlaceholderReason(real) {
		t.Errorf("%q - осмысленная причина, её нельзя терять", real)
	}
}

// Формат причины по нехватке средств повторяет эталон партнёра
func TestBuildDeclineReasonInsufficientFunds(t *testing.T) {
	cardInfo := &d8corp.CardInfoData{
		CardBasicInfo: d8corp.CardBasicInfo{
			StatCode: "00",
			Lkey:     d8corp.Lkey{Pan: "5058270530003879"},
		},
		CardAccounts: []d8corp.CardAccount{{
			AccountNumber: "20216972300001176308",
			StatCode:      "00",
			AvlBal:        5.58,
		}},
	}

	got := buildDeclineReason(utils.TwoInsufficientFunds, cardInfo, true, 1000)
	want := "Insufficient funds on account #20216972300001176308. Tranx amount=1000.00 is more than " +
		"AcctEffectiveBalance=5.58 (AvailBalance=5.58, OverdraftLimit=0.00, Bonus=0.00, " +
		"TmpOverdraft=0.00, EMVOfflineHold=0.00, Protected Amount=0.00)"
	if got != want {
		t.Errorf("причина отказа:\n получено %q\n ожидалось %q", got, want)
	}

	// Статус карты по-прежнему важнее кода операции
	cardInfo.CardBasicInfo.StatCode = "12" // Card reported lost
	if got := buildDeclineReason("40", cardInfo, true, 1000); !strings.Contains(got, "card status 'Lost'") {
		t.Errorf("отказ по статусу карты потерян: %q", got)
	}
}

// Повтор операции не должен приводить к новой транзакции: отдаём прежний ответ
func TestPosDedup(t *testing.T) {
	cache := posDedup{entries: make(map[string]dedupEntry)}
	req := Request{TranNumber: "A-100", PAN: "5058270530003879", TranCode: 175, Amount: 1.0}

	key := dedupKey(req)
	if key == "" {
		t.Fatal("номер операции задан - ключ обязан быть")
	}
	if _, ok := cache.get(key); ok {
		t.Error("пустой кеш не должен ничего отдавать")
	}

	resp := new(Envelope)
	resp.Body.POSRequestRp.Response.ThisTranId = "777"
	cache.put(key, resp)

	got, ok := cache.get(key)
	if !ok || got.Body.POSRequestRp.Response.ThisTranId != "777" {
		t.Errorf("повтор должен вернуть прежний ответ, получено %+v (ok=%v)", got, ok)
	}

	// Та же сумма и карта, но другой номер операции - это новый платёж
	other := req
	other.TranNumber = "A-101"
	if _, ok := cache.get(dedupKey(other)); ok {
		t.Error("другой номер операции - другая транзакция")
	}

	// Тот же номер, но другая сумма - тоже не дубль
	changed := req
	changed.Amount = 2.0
	if _, ok := cache.get(dedupKey(changed)); ok {
		t.Error("та же операция с другой суммой не считается повтором")
	}
}

// Без номера операции дедупликация не работает - поведение прежнее
func TestPosDedupWithoutTranNumber(t *testing.T) {
	cache := posDedup{entries: make(map[string]dedupEntry)}
	key := dedupKey(Request{PAN: "5058270530003879", TranCode: 175, Amount: 1.0})
	if key != "" {
		t.Fatalf("без TranNumber ключа быть не должно, получено %q", key)
	}

	resp := new(Envelope)
	cache.put(key, resp)
	if _, ok := cache.get(key); ok {
		t.Error("пустой ключ не должен попадать в кеш")
	}
}

// Протухшие записи не отдаются и вычищаются
func TestPosDedupExpiry(t *testing.T) {
	cache := posDedup{entries: make(map[string]dedupEntry)}
	key := "stale"
	cache.entries[key] = dedupEntry{resp: Envelope{}, at: time.Now().Add(-dedupTTL - time.Second)}

	if _, ok := cache.get(key); ok {
		t.Error("запись старше окна повтора не должна отдаваться")
	}
	if _, exists := cache.entries[key]; exists {
		t.Error("протухшая запись должна удаляться")
	}
}
