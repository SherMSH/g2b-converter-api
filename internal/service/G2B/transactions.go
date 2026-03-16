package service

import (
	"converterapi/internal/config"
	"converterapi/internal/models"
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
)

// InitiateTransaction инициирует транзакцию и сгенерирует новый номер ссылки на транзакцию электронной коммерции (ecTxRefno).
// Нужен до вызова AuthorizeTransaction и ReverseTransaction.
func InitiateTransaction() (*string, error) {
	var (
		resp    d8corp.CommonResp
		ectxNum d8corp.InitTransactionResp
	)
	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/kernel/1.0/initiateTransaction", nil, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b initiateTransaction request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b initiateTransaction resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b initiateTransaction RESP marshaling err: %v", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Data, &ectxNum)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b initiateTransaction DATA marshaling err: %v", err)
		return nil, err
	}
	if ectxNum.ECTxRefNo == nil {
		logger.Errorf("[SERVICE] D8 G2b initiateTransaction err: empty ECTxRefNo")
		return nil, err
	}
	return ectxNum.ECTxRefNo, nil
}

// AuthorizeTransaction сгенерирует онлайн-транзакцию и отправит на обработку.
// На основании PAN-номера карты транзакция будет направлена ​​в банк-эмитент или внешнюю сеть (MC/VISA) для авторизации.
// Ответ в режиме реального времени будет возвращен в полях.actionCode и status.rspCode.
// ВАЖНО: В случае, если не получаете ответное сообщение, вы обязаны проверить статус транзакции,
// вызвав "GetTransactionStatus".
// При определенных обстоятельствах возможна ситуация, когда ответное сообщение теряется, но система обработки не знает об этом,
// и транзакция успешно завершается как в системе обработки, так и в системе эмитента карты.
func AuthorizeTransaction(input models.TrnInputIface, ecTxRefNo string) (*d8corp.CommonResp, error) {
	resp := &d8corp.CommonResp{}
	trnData := &d8corp.TrnData{}
	req := d8corp.AuthTxReq{
		CardKey: d8corp.CardKey{
			Pan:        input.GetPan(),
			ExpiryDate: input.GetExpDate(),
		},
		EcTxRefno:          ecTxRefNo,
		TxnType:            "TRANSF_C2A",
		TxnAmount:          input.GetAmount(),
		TxnCurrency:        input.GetCurrency(),
		TermCode:           "TRM00001",
		CrdacptID:          "MRC00001",
		CrdacptBus:         5999, //Card Acceptor Business Code
		MessageFunction:    0,    //0-Request, 2-Advice
		RecipientAccount:   input.GetRecipientAcc(),
		DestinationAccType: "00",
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b authorizeTransaction REQ marshaling err: %v", err)
		return nil, fmt.Errorf("[SERVICE] D8 G2b authorizeTransaction REQ marshaling err")
	}

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/kernel/1.0/authorizeTransaction", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b authorizeTransaction request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b authorizeTransaction resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b authorizeTransaction RESP marshaling err: %v", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Data, trnData)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b authorizeTransaction DATA marshaling err: %v", err)
		return nil, err
	}
	if len(trnData.TransactionResponse.EcTxRefno) == 0 {
		logger.Errorf("[SERVICE] D8 G2b authorizeTransaction err: empty response")
		return nil, err
	}
	return resp, nil
}

func GetTransactionStatus(tlId int, ecTxRefNo string) (*d8corp.CommonResp, error) {
	resp := &d8corp.CommonResp{}
	trnData := &d8corp.TxStatusData{}

	req := d8corp.ChkTxStatusReq{
		TlId: tlId,
	}
	if len(ecTxRefNo) != 0 {
		req = d8corp.ChkTxStatusReq{
			EcTxRefno: ecTxRefNo,
		}
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b authorizeTransaction REQ marshaling err: %v", err)
		return nil, fmt.Errorf("[SERVICE] D8 G2b authorizeTransaction REQ marshaling err")
	}

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/kernel/1.0/authorizeTransaction", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b authorizeTransaction request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b authorizeTransaction resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b authorizeTransaction RESP marshaling err: %v", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Data, trnData)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b authorizeTransaction DATA marshaling err: %v", err)
		return nil, err
	}
	if len(trnData.TxStatus.RspCode) == 0 {
		logger.Errorf("[SERVICE] D8 G2b authorizeTransaction err: empty response")
		return nil, err
	}

	return resp, nil
}
