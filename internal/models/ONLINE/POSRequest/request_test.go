package posrequestrq

import (
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"testing"
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
		{"конкретная причина не подменяется", "59", card("08", "00"), true, "59"},
		{"неверный PIN не подменяется", "53", card("08", "00"), true, "53"},
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
