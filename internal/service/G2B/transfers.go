package service

import (
	"converterapi/internal/models"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"fmt"
)

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

	switch req.TxnType {
	case utils.C2C:
		// cardKey - карта отправителя, recipientCardKey - карта получателя
		if recipientPan == "" {
			return fmt.Errorf("TRANSF_C2C requires recipient card")
		}
		req.RecipientCardKey = &d8corp.CardKey{Pan: recipientPan}

	case utils.C2A:
		// со счёта получателя списание не идёт: нужен его номер и тип
		if recipientAccount == "" {
			return fmt.Errorf("TRANSF_C2A requires recipient account")
		}
		req.RecipientAccount = recipientAccount
		if accType := transfer.GetDestinationAccountType(); accType != "" {
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

		if req.TxnType == utils.H2C && (req.SenderFirstName == "" || req.SenderLastName == "") {
			return fmt.Errorf("TRANSF_H2C requires sender first and last name")
		}
	}

	return nil
}

func isTransfer(txnType utils.TxnType) bool {
	switch txnType {
	case utils.C2C, utils.H2C, utils.C2A, utils.A2C:
		return true
	default:
		return false
	}
}
