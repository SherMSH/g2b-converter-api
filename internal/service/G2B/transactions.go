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
		req d8corp.InitTxReq
		// resp    d8corp.CommonResp
		ectxNum d8corp.InitTransactionResp
	)
	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b initiateTransaction REQ marshaling err: %v", err)
		return nil, fmt.Errorf("[SERVICE] D8 G2b initiateTransaction REQ marshaling err")
	}
	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/kernel/1.0/initiateTransaction", jsonReq, utils.D8TxHeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b initiateTransaction request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b initiateTransaction resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, &ectxNum)
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
func AuthorizeTransaction(input models.TrnInputIface, ecTxRefNo string) (*d8corp.TrnData, *d8corp.CommonResp, error) {
	resp := &d8corp.CommonResp{}
	trnData := &d8corp.TrnData{}
	logger.Infof("AuthorizeTransaction req ExpDate: %s", input.GetExpDate())

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
		// BusinessAppId:      "TBI", //TBI - Financial Institution offered Bank-Initiated P2P Money Transfer
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b authorizeTransaction REQ marshaling err: %v", err)
		return nil, nil, fmt.Errorf("[SERVICE] D8 G2b authorizeTransaction REQ marshaling err")
	}

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/kernel/1.0/authorizeTransaction", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b authorizeTransaction request sending err: %v", err)
		return nil, nil, err
	}
	logger.Infof("[SERVICE] D8 G2b authorizeTransaction resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b authorizeTransaction RESP marshaling err: %v", err)
		return nil, nil, err
	}
	err = json.Unmarshal(resp.Data, trnData)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b authorizeTransaction DATA marshaling err: %v", err)
		return nil, nil, err
	}
	if len(trnData.TransactionResponse.EcTxRefno) == 0 {
		logger.Errorf("[SERVICE] D8 G2b authorizeTransaction err: empty response")
		return nil, nil, err
	}
	return trnData, resp, nil
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
		logger.Errorf("[SERVICE] D8 G2b GetTransactionStatus REQ marshaling err: %v", err)
		return nil, fmt.Errorf("[SERVICE] D8 G2b GetTransactionStatus REQ marshaling err")
	}

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/kernel/1.0/authorizeTransaction", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b GetTransactionStatus request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b GetTransactionStatus resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b GetTransactionStatus RESP marshaling err: %v", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Data, trnData)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b GetTransactionStatus DATA marshaling err: %v", err)
		return nil, err
	}
	if len(trnData.TxStatus.RspCode) == 0 {
		logger.Errorf("[SERVICE] D8 G2b GetTransactionStatus err: empty response")
		return nil, err
	}

	return resp, nil
}

// ReverseTransaction сгенерирует операцию отмены и поместит ее в фоновую очередь для выполнения
// ВАЖНО: Перед вызовом reverseTransaction необходимо вызвать сервис xkernel/initiateTransaction и получить новый ecTxRefNo
// Номер ссылки исходной транзакции в поле originalEcTxRefno
func ReverseTransaction(input models.TrnInputIface, ecTxRefNo, originalEcTxRefno string) (*d8corp.CommonResp, error) {
	resp := &d8corp.CommonResp{}
	trnData := &d8corp.TxResponseData{}

	req := d8corp.ReverceTxReq{
		EcTxRefno:         ecTxRefNo,
		OriginalEcTxRefno: originalEcTxRefno,
		ReasonCode:        4000,
		ReversalAmount:    input.GetAmount(),
		TxnCurrency:       input.GetCurrency(),
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b reverseTransaction REQ marshaling err: %v", err)
		return nil, fmt.Errorf("[SERVICE] D8 G2b reverseTransaction REQ marshaling err")
	}

	data, status, err := utils.SendRequest("POST", config.Config.Processing.Address+"/xapi/kernel/1.0/authorizeTransaction", jsonReq, utils.D8HeadersMap)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b reverseTransaction request sending err: %v", err)
		return nil, err
	}
	logger.Infof("[SERVICE] D8 G2b reverseTransaction resp status: %v, body: %v", status, string(data))

	err = json.Unmarshal(data, resp)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b reverseTransaction RESP marshaling err: %v", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Data, trnData)
	if err != nil {
		logger.Errorf("[SERVICE] D8 G2b reverseTransaction DATA marshaling err: %v", err)
		return nil, err
	}
	if len(trnData.TxResponse.RspCode) == 0 {
		logger.Errorf("[SERVICE] D8 G2b reverseTransaction err: empty response")
		return nil, err
	}

	return resp, nil
}

func GetTransactionInfoG2b() error {
	return nil
}
