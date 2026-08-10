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
	businessAppId    string
	cvv2             string
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
func (i trnInput) GetBusinessAppId() string          { return i.businessAppId }
func (i trnInput) GetCvv2() string                   { return i.cvv2 }

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

func TestFillTransferFieldsBusinessAppId(t *testing.T) {
	tests := []struct {
		name    string
		txnType utils.TxnType
		in      trnInput
		want    string
	}{
		{
			name:    "C2C по умолчанию person-to-person",
			txnType: utils.C2C,
			in:      trnInput{txnType: utils.C2C, recipientPan: recipientPan},
			want:    "TPP",
		},
		{
			name:    "C2A по умолчанию между счетами",
			txnType: utils.C2A,
			in:      trnInput{txnType: utils.C2A, recipientAccount: someAccount},
			want:    "TAA",
		},
		{
			name:    "значение партнёра важнее умолчания",
			txnType: utils.C2C,
			in:      trnInput{txnType: utils.C2C, recipientPan: recipientPan, businessAppId: "TCP"},
			want:    "TCP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := d8corp.AuthTxReq{TxnType: tt.txnType}
			if err := fillTransferFields(&req, tt.in); err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}
			if req.BusinessAppId != tt.want {
				t.Errorf("businessAppId = %q, want %q", req.BusinessAppId, tt.want)
			}
		})
	}
}

// Спецификация запрещает тип счёта получателя "03" для TRANSF_C2A
func TestFillTransferFieldsC2ARejectsCardAccountType(t *testing.T) {
	req := d8corp.AuthTxReq{TxnType: utils.C2A}
	in := trnInput{txnType: utils.C2A, recipientAccount: someAccount, destAccType: cardAccountType}

	if err := fillTransferFields(&req, in); err == nil {
		t.Error("ожидалась ошибка на тип счёта 03 для C2A")
	}
}

// Для переводов на карту тип счёта получателя обязан быть "03"
func TestFillTransferFieldsCardDestinationType(t *testing.T) {
	req := d8corp.AuthTxReq{TxnType: utils.C2C, DestinationAccType: "00"}
	in := trnInput{txnType: utils.C2C, recipientPan: recipientPan}

	if err := fillTransferFields(&req, in); err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if req.DestinationAccType != cardAccountType {
		t.Errorf("destinationAccountType = %q, want %q", req.DestinationAccType, cardAccountType)
	}
}

// Проверка счёта обязана уходить с нулевой суммой, а CVV2 - доходить до D8:
// без него процессинг не проверит код безопасности.
func TestBuildAuthTxReqAccver(t *testing.T) {
	in := trnInput{txnType: utils.Accver, pan: senderPan, cvv2: "123"}

	req, err := BuildAuthTxReq(in, "XAPI/ref")
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if req.TxnAmount != 0 {
		t.Errorf("сумма проверки счёта = %v, want 0", req.TxnAmount)
	}
	if req.Cvv2 != "123" {
		t.Errorf("CVV2 не попал в запрос: %q", req.Cvv2)
	}
}

func TestBuildAuthTxReqPassesCvv2ForSales(t *testing.T) {
	in := trnInput{txnType: utils.Sales, pan: senderPan, cvv2: "456"}

	req, err := BuildAuthTxReq(in, "XAPI/ref")
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if req.Cvv2 != "456" {
		t.Errorf("CVV2 не попал в запрос покупки: %q", req.Cvv2)
	}
	if req.TxnAmount == 0 {
		t.Error("сумма покупки не должна обнуляться")
	}
}
