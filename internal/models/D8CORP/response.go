package d8corp

import "encoding/json"

type CommonResp struct {
	Data   json.RawMessage `json:"data"`
	Paging Paging          `json:"paging"`
	Status Status          `json:"status"`
}

type Paging struct {
	NextPagePresent bool `json:"nextPagePresent"`
	TotalCount      int  `json:"totalCount"`
}

type Status struct {
	Code        string `json:"code"`
	RspCode     string `json:"rspcode"`
	Message     string `json:"message"`
	ReasonClass int    `json:"reasonClass"`
	SessionId   string `json:"sessionId"`
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
	CardBasicInfo        CardBasicInfo         `json:"cardBasicInfo"`
	CardAccounts         []CardAccount         `json:"cardAccounts"`
	CardLimits           []CardLimit           `json:"cardLimits"`
	CardAccountLimits    []CardAccountLimit    `json:"cardAccountLimits"`
	CardAuthRestrictions []CardAuthRestriction `json:"cardAuthRestrictions"`
	CardTransactions     []CardTransaction     `json:"cardTransactions"`
	CardNotifications    []CardNotification    `json:"cardNotifications"`
}

type CardBasicInfo struct {
	Lkey             Lkey   `json:"lkey"`
	ExpiryDate       string `json:"expiryDate"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	EmbossName       string `json:"embossName"`
	DeliveryPoint    string `json:"deliveryPoint"`
	DeliveryMethod   int    `json:"deliveryMethod"`
	CustomerCode     string `json:"customerCode"`
	Currcode         string `json:"currcode"`
	CompanyRegNumber string `json:"companyRegNumber"`
	StatCode         string `json:"statCode"`
	StatReason       string `json:"statReason"`
	StatExternalUser string `json:"statExternalUser"`
	StatRequester    string `json:"statRequester"`
	StatChangeTime   string `json:"statChangeTime"`
	EffDate          string `json:"effDate"`
	Cdproduct        string `json:"cdproduct"`
	Svccode          string `json:"svccode"`
	Pvki             string `json:"pvki"`
	Pvv              string `json:"pvv"`
	Dcec             int    `json:"dcec"`
	Cvv2Type         int    `json:"cvv2Type"`
	ProductType      int    `json:"productType"`
	HceCardRef       int    `json:"hceCardRef"`
	Cvv2Tries        int    `json:"cvv2Tries"`
	Title            string `json:"title"`
}

type CardLimit struct {
	LimitCode        string  `json:"limitCode"`
	LimitType        string  `json:"limitType"`
	Priority         int     `json:"priority"`
	Objectid         string  `json:"objectid"`
	ActFrom          string  `json:"actFrom"`
	ActTo            string  `json:"actTo"`
	CurrCode         string  `json:"currCode"`
	CompanyRegNumber string  `json:"companyRegNumber"`
	TxnAmtMin        float64 `json:"txnAmtMin"`
	TxnAmtMax        float64 `json:"txnAmtMax"`
	CycAmtMax        float64 `json:"cycAmtMax"`
	CycNumMax        int     `json:"cycNumMax"`
	CycCntAmt        float64 `json:"cycCntAmt"`
	CycCntNum        int     `json:"cycCntNum"`
}

type CardAccountLimit struct {
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
	AccountNumber string  `json:"accountNumber"`
	Currency      string  `json:"currency"`
	TypeCode      string  `json:"typeCode"`
	TypeCodeDescr string  `json:"typeCodeDescr"`
	CustomerCode  string  `json:"customerCode"`
	AvlBal        float64 `json:"avlBal"`
	BlkAmt        float64 `json:"blkAmt"`
	BlkAmtCr      float64 `json:"blkAmtCr"`
	BlkAmtDr      float64 `json:"blkAmtDr"`
	OpenDate      string  `json:"openDate"`
	LastUsage     string  `json:"lastUsage"`
	StatCode      string  `json:"statCode"`
	Priority      int     `json:"priority"`
	Crlimit       float64 `json:"crlimit"`
	Balincr       float64 `json:"balincr"`
	Balincrexp    string  `json:"balincrexp"`
	TstampBalance string  `json:"tstampBalance"`
}

type CardTransaction struct {
	TlId               int     `json:"tlId"`
	EcTxRefno          string  `json:"ecTxRefno"`
	BusDate            string  `json:"busDate"`
	Stan               int     `json:"stan"`
	Rrn                string  `json:"rrn"`
	DateExp            string  `json:"dateExp"`
	TermType           string  `json:"termType"`
	Lkey               Lkey    `json:"lkey"`
	CrdacptID          string  `json:"crdacptID"`
	CrdacptBus         int     `json:"crdacptBus"`
	TermCode           string  `json:"termCode"`
	Aiid               string  `json:"aiid"`
	FnCode             int     `json:"fnCode"`
	Msgcls             string  `json:"msgcls"`
	Msgfn              string  `json:"msgfn"`
	Txnsrc             string  `json:"txnsrc"`
	Txncode            int     `json:"txncode"`
	TxStatus           int     `json:"txStatus"`
	TxnAmount          float64 `json:"txnAmount"`
	TxnCurrency        string  `json:"txnCurrency"`
	ActionCode         string  `json:"actionCode"`
	RspCode            string  `json:"rspCode"`
	ReasonCode         int     `json:"reasonCode"`
	Fiid               string  `json:"fiid"`
	Riid               string  `json:"riid"`
	CrdactplocName     string  `json:"crdactplocName"`
	CrdactplocCity     string  `json:"crdactplocCity"`
	CrdactplocCountry  string  `json:"crdactplocCountry"`
	CrdactplocPostcode string  `json:"crdactplocPostcode"`
	Tstamp_insert      string  `json:"tstamp_insert"`
	When_created       string  `json:"when_created"`
	AccountNumber      string  `json:"accountNumber"`
	AvlBal             float64 `json:"avlBal"`
	BlkAmt             float64 `json:"blkAmt"`
	Crlimit            float64 `json:"crlimit"`
	Balincr            float64 `json:"balincr"`
	Balincrexp         string  `json:"balincrexp"`
	Amtbill            float64 `json:"amtbill"`
	Amtbillcb          float64 `json:"amtbillcb"`
	Curbill            string  `json:"curbill"`
	Ratebill           float64 `json:"ratebill"`
}

type CardNotification struct {
	NotificationTarget string `json:"notificationTarget"`
	ServiceType        string `json:"serviceType"`
}

type CVVData struct {
	Lkey          Lkey   `json:"lkey"`
	CVV2          string `json:"cvv2"`
	CVV2Type      int    `json:"cvv2Type"`
	CVV2Expiry    string `json:"cvv2Expiry"`
	EncryptedData string `json:"encryptedData"`
}

type Transaction struct {
	Details TransactionDetails `json:"transaction"`
}
type TransactionDetails struct {
	TlId                    int                     `json:"tlId"`
	EcTxRefno               string                  `json:"ecTxRefno"`
	Stan                    int                     `json:"stan"`
	Stanorg                 int                     `json:"stanorg"`
	Rrn                     string                  `json:"rrn"`
	Termtype                string                  `json:"termtype"`
	Issrtcode               string                  `json:"issrtcode"`
	Acqrtcode               string                  `json:"acqrtcode"`
	Lkey                    Lkey                    `json:"lkey"`
	DateExp                 string                  `json:"dateExp"`
	CrdacptBus              int                     `json:"crdacptBus"`
	CrdacptID               string                  `json:"crdacptID"`
	TermCode                string                  `json:"termCode"`
	Aiid                    string                  `json:"aiid"`
	DateLocal               string                  `json:"dateLocal"`
	TimeLocal               string                  `json:"timeLocal"`
	Msgcls                  string                  `json:"msgcls"`
	Msgfn                   string                  `json:"msgfn"`
	Txnsrc                  string                  `json:"txnsrc"`
	FnCode                  int                     `json:"fnCode"`
	TxnCode                 int                     `json:"txnCode"`
	TxStatus                int                     `json:"txStatus"`
	Rspsrc                  string                  `json:"rspsrc"`
	TxnAmount               float64                 `json:"txnAmount"`
	TxnCurrency             string                  `json:"txnCurrency"`
	ActionCode              string                  `json:"actionCode"`
	RspCode                 string                  `json:"rspCode"`
	Aprvlcode               string                  `json:"aprvlcode"`
	Dateset                 string                  `json:"dateset"`
	Busdate                 string                  `json:"busdate"`
	Afe                     string                  `json:"afe"`
	Ife                     string                  `json:"ife"`
	Curbill                 string                  `json:"curbill"`
	Amtbill                 float64                 `json:"amtbill"`
	Ratebill                float64                 `json:"ratebill"`
	Amtbillbal              float64                 `json:"amtbillbal"`
	Scheme                  string                  `json:"scheme"`
	CrdacptlocName          string                  `json:"crdacptlocName"`
	CrdacptlocStreet        string                  `json:"crdacptlocStreet"`
	CrdacptlocCity          string                  `json:"crdacptlocCity"`
	CrdacptlocPostcode      string                  `json:"crdacptlocPostcode"`
	CrdacptlocRegion        string                  `json:"crdacptlocRegion"`
	CrdacptlocCountry       string                  `json:"crdacptlocCountry"`
	Merchcountryorg         string                  `json:"merchcountryorg"`
	Bingrp                  string                  `json:"bingrp"`
	MerchantPhysicalAddress MerchantPhysicalAddress `json:"merchantPhysicalAddress"`
	TerminalPhysicalAddress TerminalPhysicalAddress `json:"terminalPhysicalAddress"`
	DestinationAccountType  string                  `json:"destinationAccountType"`
}

type MerchantPhysicalAddress struct {
	LocationName  string `json:"locationName"`
	PostalCode    string `json:"postalCode"`
	Country       string `json:"country"`
	ContactPhone1 string `json:"contactPhone1"`
	Url           string `json:"url"`
}

type TerminalPhysicalAddress struct {
	LocationName   string `json:"locationName"`
	Address1       string `json:"address1"`
	City           string `json:"city"`
	Region         string `json:"region"`
	PostalCode     string `json:"postalCode"`
	Country        string `json:"country"`
	ContactPhone1  string `json:"contactPhone1"`
	Url            string `json:"url"`
	AdditionalInfo string `json:"additionalInfo"`
}
