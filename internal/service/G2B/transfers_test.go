package service

import (
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"testing"
)

// trnInput - минимальная реализация models.TrnInputIface и
// models.TrnTransferIface для проверки сборки запроса без обращения к D8.
type trnInput struct {
	txnType          utils.TxnType
	pan              string
	recipientPan     string
	recipientAccount string
	senderAccount    string
	destAccType      string
}

func (i trnInput) GetTxnType() utils.TxnType         { return i.txnType }
func (i trnInput) GetPan() string                    { return i.pan }
func (i trnInput) GetMBR() string                    { return "" }
func (i trnInput) GetExpDate() string                { return "3004" }
func (i trnInput) GetAccount() string                { return i.senderAccount }
func (i trnInput) GetAmount() float64                { return 1 }
func (i trnInput) GetCurrency() string               { return utils.TJSCurrency }
func (i trnInput) GetTerminal() string               { return "J527393" }
func (i trnInput) GetAcceptorID() string             { return "Test001" }
func (i trnInput) GetRecipientPan() string           { return i.recipientPan }
func (i trnInput) GetRecipientAccount() string       { return i.recipientAccount }
func (i trnInput) GetSenderAccount() string          { return i.senderAccount }
func (i trnInput) GetDestinationAccountType() string { return i.destAccType }

const (
	senderPan    = "5058270530003879"
	recipientPan = "5058270530000073"
	someAccount  = "20216972300001176308"
)

func TestFillTransferFieldsNonTransferUntouched(t *testing.T) {
	req := d8corp.AuthTxReq{TxnType: utils.Sales, CardKey: d8corp.CardKey{Pan: senderPan}}
	if err := fillTransferFields(&req, trnInput{txnType: utils.Sales}); err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if req.RecipientCardKey != nil || req.SenderAccount != "" {
		t.Errorf("покупка не должна получать реквизиты перевода: %+v", req)
	}
}

func TestFillTransferFieldsC2C(t *testing.T) {
	req := d8corp.AuthTxReq{TxnType: utils.C2C, CardKey: d8corp.CardKey{Pan: senderPan}}
	in := trnInput{txnType: utils.C2C, pan: senderPan, recipientPan: recipientPan}

	if err := fillTransferFields(&req, in); err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if req.RecipientCardKey == nil || req.RecipientCardKey.Pan != recipientPan {
		t.Errorf("карта получателя не проставлена: %+v", req.RecipientCardKey)
	}
	if req.CardKey.Pan != senderPan {
		t.Errorf("карта отправителя изменилась: %s", req.CardKey.Pan)
	}
}

func TestFillTransferFieldsC2A(t *testing.T) {
	req := d8corp.AuthTxReq{TxnType: utils.C2A, DestinationAccType: "00"}
	in := trnInput{txnType: utils.C2A, recipientAccount: someAccount, destAccType: "11"}

	if err := fillTransferFields(&req, in); err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if req.RecipientAccount != someAccount {
		t.Errorf("счёт получателя не проставлен: %q", req.RecipientAccount)
	}
	if req.DestinationAccType != "11" {
		t.Errorf("тип счёта получателя не проставлен: %q", req.DestinationAccType)
	}
}

// По спецификации 8.5 для A2C cardKey должен совпадать с recipientCardKey:
// карта здесь - получатель, источник задан счётом отправителя.
func TestFillTransferFieldsA2CSwapsCardKey(t *testing.T) {
	req := d8corp.AuthTxReq{TxnType: utils.A2C, CardKey: d8corp.CardKey{Pan: senderPan}}
	in := trnInput{txnType: utils.A2C, recipientPan: recipientPan, senderAccount: someAccount}

	if err := fillTransferFields(&req, in); err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if req.CardKey.Pan != recipientPan {
		t.Errorf("cardKey должен стать картой получателя, а стал %q", req.CardKey.Pan)
	}
	if req.RecipientCardKey == nil || req.RecipientCardKey.Pan != recipientPan {
		t.Errorf("recipientCardKey не совпадает с cardKey: %+v", req.RecipientCardKey)
	}
	if req.SenderAccount != someAccount {
		t.Errorf("счёт отправителя не проставлен: %q", req.SenderAccount)
	}
}

func TestFillTransferFieldsMissingFields(t *testing.T) {
	tests := []struct {
		name string
		req  d8corp.AuthTxReq
		in   trnInput
	}{
		{
			name: "C2C без карты получателя",
			req:  d8corp.AuthTxReq{TxnType: utils.C2C},
			in:   trnInput{txnType: utils.C2C},
		},
		{
			name: "C2A без счёта получателя",
			req:  d8corp.AuthTxReq{TxnType: utils.C2A},
			in:   trnInput{txnType: utils.C2A},
		},
		{
			name: "A2C без счёта отправителя",
			req:  d8corp.AuthTxReq{TxnType: utils.A2C},
			in:   trnInput{txnType: utils.A2C, recipientPan: recipientPan},
		},
		{
			name: "H2C без имени отправителя",
			req:  d8corp.AuthTxReq{TxnType: utils.H2C},
			in:   trnInput{txnType: utils.H2C, recipientPan: recipientPan, senderAccount: someAccount},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fillTransferFields(&tt.req, tt.in); err == nil {
				t.Errorf("ожидалась ошибка, запрос собран: %+v", tt.req)
			}
		})
	}
}
