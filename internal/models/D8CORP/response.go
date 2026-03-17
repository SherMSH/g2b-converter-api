package d8corp

import "encoding/json"

type CommonResp struct {
	Data   []byte `json:"data"`
	Paging Paging `json:"paging"`
	Status Status `json:"status"`
}

type Paging struct {
	NextPagePresent bool `json:"nextPagePresent"`
	TotalCount      int  `json:"totalCount"`
}

type Status struct {
	Code        string          `json:"code"`
	RspCode     string          `json:"rspcode"`
	Message     json.RawMessage `json:"message"`
	ReasonClass int             `json:"reasonClass"`
	SessionId   string          `json:"sessionId"`
}

type InitTransactionResp struct {
	ECTxRefNo *string `json:"ectxrefno"`
}

type TrnData struct {
	TransactionResponse TransactionResponse `json:"transactionResponse"`
	Lkey                Lkey                `json:"lkey"`
	RecipientLkey       Lkey                `json:"recipientLkey"`
}

type TransactionResponse struct {
	TlId       int    `json:"tlId"`
	EcTxRefno  string `json:"ecTxRefno"`
	Stan       int    `json:"stan"`
	Rrn        string `json:"rrn"`
	ActionCode string `json:"actionCode"`
	RspCode    string `json:"rspCode"`
}

type Lkey struct {
	MaskedPan string `json:"maskedPan"`
	Pan       string `json:"pan"`
	Sequence  int    `json:"sequence"`
	LkeyId    int    `json:"lkeyId"`
	LkeyAlias string `json:"lkeyAlias"`
}

type TxStatusData struct {
	TxStatus TxStatus `json:"transactionStatus"`
}

type TxStatus struct {
	Rrn         string `json:"rrn"`
	ActionCode  string `json:"actionCode"`
	RspCode     string `json:"rspCode"`
	TxStatusInt int    `json:"txStatus"`
}

type TxResponseData struct {
	TxResponse TxResponse `json:"transactionResponse"`
}

type TxResponse struct {
	TlId       int    `json:"tlId"`
	EcTxRefNo  string `json:"ecTxRefno"`
	Stan       int    `json:"stan"`
	ActionCode string `json:"actionCode"`
	RspCode    string `json:"rspCode"`
	TxStatus   int    `json:"txStatus"`
}
