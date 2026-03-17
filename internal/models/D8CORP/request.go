package d8corp

type Request struct {
	MDIRecords []MdiRecord `json:"mdiRecords"`
}

type MdiRecord struct {
	IssRectype           string `json:"ISS_RECTYPE"`
	IssRecaction         string `json:"ISS_RECACTION"`
	IssRecnum            int    `json:"ISS_RECNUM"`
	IssCompanyRegnr      string `json:"ISS_COMPANY_REGNR"`
	IssCompanyRegnrAcc   string `json:"ISS_COMPANY_REGNR_ACC"`
	IssImpPvki           string `json:"ISS_IMP_PVKI"`
	DbCustomerCustcode   string `json:"DB_CUSTOMER_CUSTCODE"`
	DbCdproductCdproduct string `json:"DB_CDPRODUCT_CDPRODUCT"`
	DbAccountAccnum      string `json:"DB_ACCOUNT_ACCNUM"`
	DbAccountCurrcode    string `json:"DB_ACCOUNT_CURRCODE"`
	KlLkeyAlias          string `json:"KL_LKEY_ALIAS"`
	KlLkeySeqno          string `json:"KL_LKEY_SEQNO"`
	DbCardaExpdate       string `json:"DB_CARDA_EXPDATE"`
	DbCardaCommCat       string `json:"DB_CARDA_COMM_CAT"`
	DbCardaEnroll3ds     string `json:"DB_CARDA_ENROLL3DS"`
	DbCardaLimitCat      string `json:"DB_CARDA_LIMIT_CAT"`
	DbCardEmbossname     string `json:"DB_CARD_EMBOSSNAME"`
	DbCardFirstname      string `json:"DB_CARD_FIRSTNAME"`
	DbCardLastname       string `json:"DB_CARD_LASTNAME"`
	DbCardMaidenname     string `json:"DB_CARD_MAIDENNAME"`
	DbCardDeliveryPoint  string `json:"DB_CARD_DELIVERY_POINT"`
	DbCrdaccPriority     string `json:"DB_CRDACC_PRIORITY"`
}
type InitTxReq struct {
}

type AuthTxReq struct {
	EcTxRefno          string  `json:"ecTxRefno"`         //+
	TxnType            string  `json:"txnType"`           //+
	CardKey            CardKey `json:"cardKey,omitempty"` //+
	TxnAmount          float64 `json:"txnAmount"`         //+
	TxnCurrency        string  `json:"txnCurrency"`       //+
	TermCode           string  `json:"termCode"`          //+
	CrdacptID          string  `json:"crdacptID"`         //-
	MessageFunction    int     `json:"messageFunction"`   //+
	MerchantCommission float64 `json:"merchantCommission,omitempty"`
	CrdacptBus         int     `json:"crdacptBus,omitempty"`
	MerchantName       string  `json:"merchantName,omitempty"`
	MerchantStreet     string  `json:"merchantStreet,omitempty"`
	MerchantCity       string  `json:"merchantCity,omitempty"`
	MerchantPostcode   string  `json:"merchantPostcode,omitempty"`
	MerchantCountry    string  `json:"merchantCountry,omitempty"`
	MerchantRegion     string  `json:"merchantRegion,omitempty"`
	Cvv2               string  `json:"cvv2,omitempty"`
	SenderFirstName    string  `json:"senderfirstName,omitempty"`
	SenderLastName     string  `json:"senderlastName,omitempty"`
	SenderLocCity      string  `json:"senderlocCity,omitempty"`
	SenderLocStreet    string  `json:"senderlocStreet,omitempty"`
	SenderLocCountry   string  `json:"senderlocCountry,omitempty"`
	RecipientCountry   string  `json:"recipientCountry,omitempty"`
	RecipientCity      string  `json:"recipientCity,omitempty"`
	RecipientStreet    string  `json:"recipientStreet,omitempty"`
	RecipientFirstName string  `json:"recipientfirstName,omitempty"`
	RecipientLastName  string  `json:"recipientlastName,omitempty"`
	RecipientAccount   string  `json:"recipientAccount,omitempty"`
	DestinationAccType string  `json:"destinationAccountType,omitempty"`
	SenderFundsrc      string  `json:"senderFundsrc,omitempty"`
	BusinessAppId      string  `json:"businessAppId,omitempty"`
	Eci3DS             string  `json:"eci3DS,omitempty"`
	SecurityFlag       int     `json:"securityFlag,omitempty"`
	Ucaf               string  `json:"ucaf,omitempty"`
	IssCommCode        string  `json:"issCommCode,omitempty"`
	AcqCommCode        string  `json:"acqCommCode,omitempty"`
	PurchRefNo         string  `json:"purchrefno,omitempty"`
}

type ChkTxStatusReq struct {
	EcTxRefno string `json:"ecTxRefno,omitempty"`
	TlId      int    `json:"tlId,omitempty"`
}

type CardKey struct {
	Pan        string `json:"pan"`
	ExpiryDate string `json:"expiryDate"`
}

type ReverceTxReq struct {
	EcTxRefno         string  `json:"ecTxRefno"`
	OriginalEcTxRefno string  `json:"originalEcTxRefno"`
	ReasonCode        int     `json:"reasonCode"`
	ReversalAmount    float64 `json:"reversalAmount"`
	TxnCurrency       string  `json:"txnCurrency"`
}
