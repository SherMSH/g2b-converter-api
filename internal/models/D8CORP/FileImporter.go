package d8corp

import (
	"converterapi/internal/utils"
	"encoding/json"
)

type MdiFile struct {
	MdiRecords []json.RawMessage `json:"mdiRecords"`
}

type MdiRecordDetails struct {
	IssRectype            string              `json:"ISS_RECTYPE"` // {"MERCHANT", "CARD, "ACCOUNT", "CUSTOMER", "POS", "CDRNOTIF"}
	IssRecaction          utils.Action        `json:"ISS_RECACTION"`
	IssRecnum             int                 `json:"ISS_RECNUM"`
	IssCompanyRegnr       utils.CompanyRegNum `json:"ISS_COMPANY_REGNR"`
	IssCompanyRegnrNew    utils.CompanyRegNum `json:"ISS_COMPANY_REGNR_NEW,omitempty"`
	IssCompanyRegnrAcc    string              `json:"ISS_COMPANY_REGNR_ACC,omitempty"`
	IssImpPvki            int                 `json:"ISS_IMP_PVKI,omitempty"` //PVKI used for PVV validation (should be values 0..6). Field required in case if ISS_GEN_PIN=3.
	DbCustomerCustcode    string              `json:"DB_CUSTOMER_CUSTCODE"`
	DbCustomerTypeCode    int                 `json:"DB_CUSTOMER_TYPECODE,omitempty"`
	DbCdproductCdproduct  string              `json:"DB_CDPRODUCT_CDPRODUCT"`
	DbAccountAccnum       string              `json:"DB_ACCOUNT_ACCNUM,omitempty"`
	DbAccountTypecode     string              `json:"DB_ACCOUNT_TYPECODE,omitempty"`
	DbAccountCurrcode     string              `json:"DB_ACCOUNT_CURRCODE,omitempty"`
	DbCrdaccPriority      int                 `json:"DB_CRDACC_PRIORITY,omitempty"`
	KlLkeyAlias           string              `json:"KL_LKEY_ALIAS,omitempty"` //`В случае ISS_RECTYPE<>CARD в вызове сообщения должно присутствовать либо KL_LKEY_ALIAS, либо KL_LKEY_CLR. В случае ISS_RECTYPE=CARD в вызове сообщения должно присутствовать одно из следующих полей/комбинаций полей (приоритет учитывается): 1. DB_CARDA_ID (если поле доступно в текущем типе файла) или 2. KL_LKEY_CLR + DB_CARDA_EXPDATE + необязательно: KL_LKEY_SEQNO или 3. KL_LKEY_ALIAS + DB_CARDA_EXPDATE + необязательно: KL_LKEY_SEQNO`
	KlLKeyClr             string              `json:"KL_LKEY_CLR,omitempty"`   // Payment instrument clear value
	KlLkeySeqno           int                 `json:"KL_LKEY_SEQNO,omitempty"`
	DbCardaExpdate        int                 `json:"DB_CARDA_EXPDATE,omitempty"`
	DbCardaStatcode       string              `json:"DB_CARDA_STATCODE,omitempty"`
	DbCardaCommCat        string              `json:"DB_CARDA_COMM_CAT,omitempty"`
	DbCardaFeeCat         string              `json:"DB_CARDA_FEE_CAT,omitempty"`
	DbCardaLimitCat       string              `json:"DB_CARDA_LIMIT_CAT,omitempty"`
	DbCardaAuthChkCat     string              `json:"DB_CARDA_AUTHCHK_CAT,omitempty"`
	DbCardaEnroll3ds      string              `json:"DB_CARDA_ENROLL3DS,omitempty"`
	DbCardEmbossname      string              `json:"DB_CARD_EMBOSSNAME,omitempty"`
	DbCardFirstname       string              `json:"DB_CARD_FIRSTNAME,omitempty"`
	DbCardLastname        string              `json:"DB_CARD_LASTNAME,omitempty"`
	DbCardTitle           string              `json:"DB_CARD_TITLE,omitempty"`
	DbCardFathername      string              `json:"DB_CARD_FATHERNAME,omitempty"`
	DbCardMaidenname      string              `json:"DB_CARD_MAIDENNAME,omitempty"`
	DbCardDeliveryPoint   string              `json:"DB_CARD_DELIVERY_POINT,omitempty"`
	DbCardaDesignImage    string              `json:"DB_CARDA_DESIGN_IMAGE,omitempty"`
	DbCardaOptions        string              `json:"DB_CARDA_OPTIONS"`
	DbCustomerFirstName   string              `json:"DB_CUSTOMER_FIRSTNAME,omitempty"`
	DbCustomerLastName    string              `json:"DB_CUSTOMER_LAASTNAME,omitempty"`
	DbCustomerDateBirth   string              `json:"DB_CUSTOMER_DATE_BIRTH,omitempty"` //format YYYYMMDD
	DbCustomerHomeCountry int                 `json:"DB_CUSTOMER_HOME_COUNTRY,omitempty"`

	DbCdNotifSvcTyp string `json:"DB_CDNOTIF_SVCTYP,omitempty"` //SMSGEN - Generic SMS notifications; SMSTXN - Transaction notifications by SMS
	DbCdNotifTarget string `json:"DB_CDNOTIF_TARGET,omitempty"` //mobile phone, e-mail or ect.
}

type HeaderRecord struct {
	IssRectype      string              `json:"ISS_RECTYPE"`       // Const "HEADER"
	IssRecaction    string              `json:"ISS_RECACTION"`     // Const "IMPORT"
	CFilename       string              `json:"C_FILENAME"`        // Format “G2BISS-YYYYMMSS-HHMISS.JSON”
	IssSourcesys    string              `json:"ISS_SOURCESYS"`     // const "ConverterApi"
	IssCompanyRegnr utils.CompanyRegNum `json:"ISS_COMPANY_REGNR"` // one of utils.CompanuRegNum consts
	IssTimestamp    string              `json:"ISS_TIMESTAMP"`     // YYYYMMDDHHMISSfff
}

type FooterRecord struct {
	IssRectype   string `json:"ISS_RECTYPE"`   // Const "FOOTER"
	IssRecaction string `json:"ISS_RECACTION"` // Const "IMPORT"
	CFilename    string `json:"C_FILENAME"`    // Format “G2BISS-YYYYMMSS-HHMISS.JSON”
	IssReccnt    int    `json:"ISS_RECCNT"`    // Number of Records
}
