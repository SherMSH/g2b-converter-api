package service

import (
	"converterapi/internal/models"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"fmt"
)

// cardAccountType - тип счёта получателя "карточный счёт".
// Спецификация (8.5) требует его для TRANSF_C2C, TRANSF_A2C и TRANSF_H2C и
// прямо запрещает для TRANSF_C2A.
const cardAccountType = "03"

// BuildAuthTxReq собирает запрос авторизации.
//
// Вынесен отдельно, потому что xkernel/calculateAcqCommission принимает ровно
// тот же запрос, что и xkernel/authorizeTransaction (8.8).
func BuildAuthTxReq(input models.TrnInputIface, ecTxRefNo string) (d8corp.AuthTxReq, error) {
	if input.GetTxnType() == "" {
		return d8corp.AuthTxReq{}, fmt.Errorf("unsupported transaction type")
	}

	req := d8corp.AuthTxReq{
		CardKey: d8corp.CardKey{
			Pan:        input.GetPan(),
			ExpiryDate: input.GetExpDate(),
		},
		EcTxRefno:          ecTxRefNo,
		TxnType:            input.GetTxnType(),
		TxnAmount:          input.GetAmount(),
		TxnCurrency:        input.GetCurrency(),
		TermCode:           input.GetTerminal(),
		CrdacptID:          input.GetAcceptorID(),
		CrdacptBus:         5999, //Card Acceptor Business Code
		MessageFunction:    0,    //0-Request, 2-Advice
		DestinationAccType: "00",
		Cvv2:               input.GetCvv2(),
	}

	// Проверка счёта: сумма обязана быть нулевой (8.5)
	if req.TxnType == utils.Accver {
		req.TxnAmount = 0
	}

	if err := fillTransferFields(&req, input); err != nil {
		return d8corp.AuthTxReq{}, err
	}
	return req, nil
}

// fillTransferFields дополняет запрос авторизации реквизитами перевода и
// проверяет обязательные комбинации полей (8.5 спецификации D8).
//
// Для непереводных операций ничего не делает.
func fillTransferFields(req *d8corp.AuthTxReq, input models.TrnInputIface) error {
	if !isTransfer(req.TxnType) {
		return nil
	}

	transfer, ok := input.(models.TrnTransferIface)
	if !ok {
		return fmt.Errorf("transaction type %s requires transfer details", req.TxnType)
	}

	recipientPan := transfer.GetRecipientPan()
	recipientAccount := transfer.GetRecipientAccount()
	senderAccount := transfer.GetSenderAccount()

	// Назначение перевода обязательно для всех TRANSF_*
	req.BusinessAppId = transfer.GetBusinessAppId()
	if req.BusinessAppId == "" {
		req.BusinessAppId = defaultBusinessAppId(req.TxnType)
	}

	switch req.TxnType {
	case utils.C2C:
		// cardKey - карта отправителя, recipientCardKey - карта получателя
		if recipientPan == "" {
			return fmt.Errorf("TRANSF_C2C requires recipient card")
		}
		req.RecipientCardKey = &d8corp.CardKey{Pan: recipientPan}
		req.DestinationAccType = cardAccountType

	case utils.C2A:
		// со счёта получателя списание не идёт: нужен его номер и тип
		if recipientAccount == "" {
			return fmt.Errorf("TRANSF_C2A requires recipient account")
		}
		req.RecipientAccount = recipientAccount
		if accType := transfer.GetDestinationAccountType(); accType != "" {
			if accType == cardAccountType {
				return fmt.Errorf("TRANSF_C2A does not allow destination account type %s", cardAccountType)
			}
			req.DestinationAccType = accType
		}

	case utils.A2C, utils.H2C:
		// по спецификации cardKey должен совпадать с recipientCardKey:
		// карта здесь - получатель, а источник задан счётом отправителя
		if recipientPan == "" {
			return fmt.Errorf("%s requires recipient card", req.TxnType)
		}
		if senderAccount == "" {
			return fmt.Errorf("%s requires sender account", req.TxnType)
		}
		req.CardKey = d8corp.CardKey{Pan: recipientPan}
		req.RecipientCardKey = &d8corp.CardKey{Pan: recipientPan}
		req.SenderAccount = senderAccount
		req.DestinationAccType = cardAccountType

		if req.TxnType == utils.H2C && (req.SenderFirstName == "" || req.SenderLastName == "") {
			return fmt.Errorf("TRANSF_H2C requires sender first and last name")
		}
	}

	return nil
}

// defaultBusinessAppId - назначение перевода по умолчанию (Appendix D):
// TPP - между людьми, TAA - между счетами.
func defaultBusinessAppId(txnType utils.TxnType) string {
	switch txnType {
	case utils.C2C, utils.H2C:
		return "TPP"
	case utils.C2A, utils.A2C:
		return "TAA"
	default:
		return ""
	}
}

func isTransfer(txnType utils.TxnType) bool {
	switch txnType {
	case utils.C2C, utils.H2C, utils.C2A, utils.A2C:
		return true
	default:
		return false
	}
}
