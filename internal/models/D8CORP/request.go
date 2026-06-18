package d8corp

import "converterapi/internal/utils"

type SetCardStatusReq struct {
	CardKey      CardKey `json:"cardKey"`
	NewStatCode  string  `json:"newStatCode"`
	Reason       string  `json:"reason"`
	Force        bool    `json:"force,omitempty"`
	ExternalUser string  `json:"externalUser,omitempty"`
}
type SetPinReq struct {
	CardKey        CardKey `json:"cardKey"`
	PinKeyUnderRSA string  `json:"pinKeyUnderRSA"`
	PinBlock       string  `json:"pinBlock,omitempty"`
	PinBlockType   string  `json:"pinBlockType,omitempty"`
}
type GetCardInfoReq struct {
	CardKey                 CardKey `json:"cardKey"`
	ReqCardBasicInfo        bool    `json:"reqCardBasicInfo,omitempty"`
	ReqCardAccounts         bool    `json:"reqCardAccounts,omitempty"`
	ReqCardLimits           bool    `josn:"reqCardLimits,omitempty"`
	ReqCardAccountLimits    bool    `json:"reqCardAccountLimits,omitempty"`
	ReqCardAuthRestrictions bool    `json:"reqCardAuthRestrictions,omitempty"`
	ReqCardTransactions     bool    `json:"reqCardTransactions,omitempty"`
	ReqCardNotifications    bool    `json:"reqCardNotifications,omitempty"`
	CardTransactionCount    int     `json:"cardTransactionCount,omitempty"`
}

type GetCardTrnHistoryReq struct {
	CardKey       CardKey      `json:"cardKey"`
	DateUTCFrom   string       `json:"dateUTCFrom"`
	DateUTCTo     string       `json:"dateUTCTo"`
	DateLocalFrom string       `json:"dateLocalFrom"`
	DateLocalTo   string       `json:"dateLocalTo"`
	PagingParams  PagingParams `json:"paging"`
}

type GetCVVReq struct {
	CardKey CardKey `json:"cardKey"`
	// RsaKeyBlock KeyBlock `json:"rsaKeyBlock,omitempty"`
}

type KeyBlock struct {
}

type InitTxReq struct {
}

type AuthTxReq struct {
	EcTxRefno          string        `json:"ecTxRefno"`         //+
	TxnType            utils.TxnType `json:"txnType"`           //+
	CardKey            CardKey       `json:"cardKey,omitempty"` //+
	TxnAmount          float64       `json:"txnAmount"`         //+
	TxnCurrency        string        `json:"txnCurrency"`       //+
	TermCode           string        `json:"termCode"`          //+
	CrdacptID          string        `json:"crdacptID"`         //-
	MessageFunction    int           `json:"messageFunction"`   //+
	MerchantCommission float64       `json:"merchantCommission,omitempty"`
	CrdacptBus         int           `json:"crdacptBus,omitempty"`
	MerchantName       string        `json:"merchantName,omitempty"`
	MerchantStreet     string        `json:"merchantStreet,omitempty"`
	MerchantCity       string        `json:"merchantCity,omitempty"`
	MerchantPostcode   string        `json:"merchantPostcode,omitempty"`
	MerchantCountry    string        `json:"merchantCountry,omitempty"`
	MerchantRegion     string        `json:"merchantRegion,omitempty"`
	Cvv2               string        `json:"cvv2,omitempty"`
	SenderFirstName    string        `json:"senderfirstName,omitempty"`
	SenderLastName     string        `json:"senderlastName,omitempty"`
	SenderLocCity      string        `json:"senderlocCity,omitempty"`
	SenderLocStreet    string        `json:"senderlocStreet,omitempty"`
	SenderLocCountry   string        `json:"senderlocCountry,omitempty"`
	RecipientCountry   string        `json:"recipientCountry,omitempty"`
	RecipientCity      string        `json:"recipientCity,omitempty"`
	RecipientStreet    string        `json:"recipientStreet,omitempty"`
	RecipientFirstName string        `json:"recipientfirstName,omitempty"`
	RecipientLastName  string        `json:"recipientlastName,omitempty"`
	RecipientAccount   string        `json:"recipientAccount,omitempty"`
	DestinationAccType string        `json:"destinationAccountType,omitempty"`
	SenderFundsrc      string        `json:"senderFundsrc,omitempty"`
	BusinessAppId      string        `json:"businessAppId,omitempty"`
	Eci3DS             string        `json:"eci3DS,omitempty"`
	SecurityFlag       int           `json:"securityFlag,omitempty"`
	Ucaf               string        `json:"ucaf,omitempty"`
	IssCommCode        string        `json:"issCommCode,omitempty"`
	AcqCommCode        string        `json:"acqCommCode,omitempty"`
	PurchRefNo         string        `json:"purchrefno,omitempty"`
}

type ChkTxStatusReq struct {
	EcTxRefno string `json:"ecTxRefno,omitempty"`
	TlId      int    `json:"tlId,omitempty"`
}

type CardKey struct {
	Lkey       int    `json:"lkeyId,omitempty"`
	Pan        string `json:"pan,omitempty"`
	ExpiryDate string `json:"expiryDate,omitempty"`
}

type ReverceTxReq struct {
	EcTxRefno         string  `json:"ecTxRefno"`
	OriginalEcTxRefno string  `json:"originalEcTxRefno"`
	ReasonCode        int     `json:"reasonCode"`
	ReversalAmount    float64 `json:"reversalAmount"`
	TxnCurrency       string  `json:"txnCurrency"`
}

type PagingParams struct {
	Size            int `json:"size"`
	LastRetrievedId int `json:"lastRetrievedId"`
}
