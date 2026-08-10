package posrequestrq

import (
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
