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
	ApprovalCode        string              `json:"aprvlCode"`
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

type CardInfoData struct {
	Lkey          Lkey   `json:"lkey"`
	ExpDate       string `json:"expiryDate"` //YYDD
	Title         string `json:"title"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	EmbossName    string `json:"embossName"`
	DeliveryPoint string `json:"deliveryPoint"`
	CustomerCode  string `json:"customerCode"`
	StatCode      string `json:"statCode"`
	// StatReason	string
	EffDate              string                `json:"effDate"`
	SvcCode              string                `json:"svccode"`
	PVV                  string                `json:"pvv"`
	InvalidPinTries      string                `json:"invalidPinTries"`
	CommissionCategory   string                `json:"comissionCategory"`
	FeeCategory          string                `json:"feeCategory"`
	LimitCategory        string                `json:"limitCategory"`
	CardAuthRestrictions []CardAuthRestriction `json:"cardAuthRestrictions"`
	CardAccounts         []CardAccount         `json:"cardAccounts"`
	CardLimits           []byte                `json:"cardLimits"`
	CardAccountLimits    []byte                `json:"cardAccountLimits"`
}

type CardAuthRestriction struct {
	TlId               int     `json:"tlId"`
	ExTxRefNo          string  `json:"ecTxRefno"`
	Stan               int     `json:"stan"`
	RRN                string  `json:"rrn"`
	DateExp            string  `json:"dateExp"`
	LkeyId             int     `json:"lkeyId"`
	CrdacptID          string  `json:"crdacptID"`
	CrdacptBus         int     `json:"crdacptBus"`
	TempCode           string  `json:"termCode"`
	Aiid               int     `json:"aiid"`
	FnCode             int     `json:"fnCode"`
	Msgcls             string  `json:"msgcls"`
	Txncode            int     `json:"txncode"`
	TxStatus           int     `json:"txStatus"`
	TxnAmount          float64 `json:"txnAmount"`
	TxnCurrency        string  `json:"txnCurrency"`
	ActionCode         string  `json:"actionCode"`
	RspCode            string  `json:"rspCode"`
	CrdactplocName     string  `json:"crdactplocName"`
	CrdactplocCity     string  `json:"crdactplocCity"`
	CrdactplocCountry  string  `json:"crdactplocCountry"`
	CrdactplocPostcode string  `json:"crdactplocPostcode"`
	Tstamp_insert      string  `json:"tstamp_insert"`
	When_created       string  `json:"when_created"`
}

type CardAccount struct {
	AcctNo string `json:"acctNo"`
}
