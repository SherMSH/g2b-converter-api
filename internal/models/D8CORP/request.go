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
	PinBlockType   int16   `json:"pinBlockType,omitempty"`
}
type GetCardInfoReq struct {
	CardKey                 CardKey `json:"cardKey"`
	ReqCardBasicInfo        bool    `json:"reqCardBasicInfo,omitempty"`
	ReqCardAccounts         bool    `json:"reqCardAccounts,omitempty"`
	ReqCardLimits           bool    `json:"reqCardLimits,omitempty"`
	ReqCardAccountLimits    bool    `json:"reqCardAccountLimits,omitempty"`
	ReqCardAuthRestrictions bool    `json:"reqCardAuthRestrictions,omitempty"`
	ReqCardTransactions     bool    `json:"reqCardTransactions,omitempty"`
	ReqCardNotifications    bool    `json:"reqCardNotifications,omitempty"`
	CardTransactionCount    int     `json:"cardTransactionCount,omitempty"`
}

type GetCardTrnHistoryReq struct {
	CardKey       CardKey      `json:"cardKey"`
	DateUTCFrom   string       `json:"dateInsertUTCFrom"`
	DateUTCTo     string       `json:"dateInsertUTCTo"`
	DateLocalFrom string       `json:"dateInsertLocalFrom"`
	DateLocalTo   string       `json:"dateInsertLocalTo"`
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
	RecipientCardKey   *CardKey      `json:"recipientCardKey,omitempty"`
	SenderAccount      string        `json:"senderAccount,omitempty"`
	SourceAccount      string        `json:"sourceAccount,omitempty"`
	SourceAccountType  string        `json:"sourceAccountType,omitempty"`
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

// AccountKey - ключ счёта (5.3.6 спецификации D8).
// В запросе должен присутствовать ровно один набор: либо Id, либо
// AccountNumber + Currency.
type AccountKey struct {
	Id            int    `json:"id,omitempty"`
	AccountNumber string `json:"accountNumber,omitempty"`
	Currency      string `json:"currency,omitempty"`
}

// PagingAdvanced - постраничная навигация вида size+page (4.x спецификации).
// Не может использоваться вместе с PagingParams.
type PagingAdvanced struct {
	Size int `json:"size"`
	Page int `json:"page"`
}

// Типы журнала выписки - поле logType, появилось в версии 1.81 спецификации.
const (
	LogTypeAll       = 0 // все записи, включая смену статусов счёта
	LogTypeFinancial = 1 // только финансовые движения
	LogTypeBlocked   = 2 // только изменения блокированных сумм
	LogTypeFinBlk    = 3 // финансовые движения и блокировки
)

// GetAccountStatementReq - запрос выписки по счёту (7.30).
// Даты в формате YYYYMMDD.
type GetAccountStatementReq struct {
	AccountKey AccountKey `json:"accountKey"`
	DateFrom   string     `json:"dateFrom"`
	DateTo     string     `json:"dateTo"`

	// LogType без omitempty: нулевое значение здесь осмысленное - "все записи"
	LogType        int            `json:"logType"`
	PagingAdvanced PagingAdvanced `json:"pagingAdvanced"`
}

// GetAccountCardListReq - запрос карт, привязанных к счёту (7.28).
type GetAccountCardListReq struct {
	AccountKey     AccountKey     `json:"accountKey"`
	PagingAdvanced PagingAdvanced `json:"pagingAdvanced"`
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
