package utils

type RqBodyType string

const (
	Unknown               RqBodyType = ""
	AcctCreditRq          RqBodyType = "AcctCreditRq"
	AcctDebitRq           RqBodyType = "AcctDebitRq"
	AddCMSAbonentRq       RqBodyType = "AddCMSAbonentRq"
	AddPersonCMSAbonentRq RqBodyType = "AddPersonCMSAbonentRq"

	// CardCredit; CardDebit
	POSRequestRq RqBodyType = "POSRequestRq"
	//
	ChangeCMSAbonentRq RqBodyType = "ChangeCMSAbonentRq"
	GetAcctInfoRq      RqBodyType = "GetAcctInfoRq"
	GetAcctStatementRq RqBodyType = "GetAcctStatementRq"
	GetCardInfoRq      RqBodyType = "GetCardInfoRq"
	GetCardStatementRq RqBodyType = "GetCardStatementRq"
	GetCVVRq           RqBodyType = "GetCVVRq"
	GetPersonInfoRq    RqBodyType = "GetPersonInfoRq"
	GetTransInfoRq     RqBodyType = "GetTransInfoRq"
	InitSessionRq      RqBodyType = "InitSessionRq"

	// RelinkPreIssuedCards
	UpdateCard2AcctLinkRq RqBodyType = "UpdateCard2AcctLinkRq"
	DeleteCard2AcctLinkRq RqBodyType = "DeleteCard2AcctLinkRq"
	SetCardPersonRq       RqBodyType = "SetCardPersonRq"
	// SetCardStatusRq       RqBodyType = "SetCardStatusRq"
	//
	ReverseTransactionRq     RqBodyType = "ReverseTransactionRq"
	RemoveCMSAbonentRq       RqBodyType = "RemoveCMSAbonentRq"
	RemovePersonCMSAbonentRq RqBodyType = "RemovePersonCMSAbonentRq"
	ResetBadPINTriesRq       RqBodyType = "ResetBadPINTriesRq"
	SetCardStatusRq          RqBodyType = "SetCardStatusRq"
	SetAcctStatusRq          RqBodyType = "SetAcctStatusRq"

	// Проверка PIN. Партнёр пишет тег и как VerifyPINRq, и как VerifyPinRq -
	// принимаем оба написания, чтобы не ловить "Unknown XML body" на регистре
	VerifyPINRq RqBodyType = "VerifyPINRq"
	VerifyPinRq RqBodyType = "VerifyPinRq"

	// Установка PIN из мобильного приложения. Обрабатывается только в
	// мобильном канале /g2b/mobile/v1/d8convert
	SetDynamicPVVPINOffsetRq RqBodyType = "SetDynamicPVV_PINOffsetRq"

	// Смена PIN. Партнёр пишет тег в двух вариантах - принимаем оба
	ChangePINRq RqBodyType = "ChangePINRq"
	ChangePIN   RqBodyType = "ChangePIN"

	UpdatePersonRq RqBodyType = "UpdatePersonRq"
)

var BodyTypes = []RqBodyType{
	AcctCreditRq,
	AcctDebitRq,
	AddCMSAbonentRq,
	AddPersonCMSAbonentRq,
	POSRequestRq,
	ChangeCMSAbonentRq,
	GetAcctInfoRq,
	GetAcctStatementRq,
	GetCardInfoRq,
	GetCardStatementRq,
	GetCVVRq,
	GetPersonInfoRq,
	GetTransInfoRq,
	InitSessionRq,
	UpdateCard2AcctLinkRq,
	DeleteCard2AcctLinkRq,
	SetCardPersonRq,
	ReverseTransactionRq,
	RemoveCMSAbonentRq,
	RemovePersonCMSAbonentRq,
	ResetBadPINTriesRq,
	SetCardStatusRq,
	SetAcctStatusRq,
	VerifyPINRq,
	VerifyPinRq,
	SetDynamicPVVPINOffsetRq,
	ChangePINRq,
	ChangePIN,
	UpdatePersonRq,
}

type OfflineReqType string

const (
	CreateCardsOut                          OfflineReqType = "CreateCardsOut*.xml"
	CreateCustomerAndAccount                OfflineReqType = "CreateCustomerAndAccount*.xml"
	CreateOrganizations                     OfflineReqType = "CreateOrganizations*.xml"
	CreatePreIssuedCards                    OfflineReqType = "CreatePreIssuedCards*.xml"
	CreateStatusActivationsOut              OfflineReqType = "CreateStatusActivationsOut*.xml"
	ReissueCardsOut                         OfflineReqType = "ReissueCardsOut*.xml"
	RelinkPreIssuedCardsOut                 OfflineReqType = "RelinkPreIssuedCardsOut*.xml"
	RelinkPreIssuedCardStatusActivationsOut OfflineReqType = "RelinkPreIssuedCardStatusActivationsOut*.xml"
)

var OfflineReqTypes = []OfflineReqType{
	CreateCardsOut,
	CreateCustomerAndAccount,
	CreateOrganizations,
	CreatePreIssuedCards,
	CreateStatusActivationsOut,
	ReissueCardsOut,
	RelinkPreIssuedCardsOut,
	RelinkPreIssuedCardStatusActivationsOut,
}

var D8HeadersMap = make(map[string]string)
var D8TxHeadersMap = make(map[string]string)

type CompanyRegNum string

const (
	Dummy  CompanyRegNum = "Default"
	Arvand CompanyRegNum = "ARV"
	Humo   CompanyRegNum = "1111"
)

type Action string

const (
	Add    = "ADD"
	Update = "UPDATE"
	Delete = "DELETE"
	Set    = "SET" // для установки адресов
)

type D8RspCode string

const (
	D8Approved         D8RspCode = "00"
	D8HonourWithID     D8RspCode = "01"
	D8PartialAmnt      D8RspCode = "02"
	ICCOfflineApproved D8RspCode = "80"
	ICCUnableGoOnline  D8RspCode = "81"
	ICCAAR             D8RspCode = "82"
	NotDeclined        D8RspCode = "85"
)

type D8TxStatus string

// Статусы транзакции по Appendix B спецификации D8.
//
// Значения 2-5 раньше были расставлены неверно: 2 считался «завершена»,
// 3 - «отклонена», хотя по спецификации 2 - это отказ, а 3 - одобрение.
// Из-за этого одобренные операции помечались ошибкой.
const (
	InProgress        D8TxStatus = "1" // получена в онлайне, статус не финальный
	Rejected          D8TxStatus = "2" // завершена, отклонена
	Approved          D8TxStatus = "3" // завершена, одобрена
	ApprovedGotRev    D8TxStatus = "4" // одобрена, получен полный реверс
	ApprovedGotPrtRev D8TxStatus = "5" // одобрена, получен частичный реверс

	AdviceLogNotProcessed D8TxStatus = "6"
	AdviceLogRejected     D8TxStatus = "7"
	AdviceLogApproved     D8TxStatus = "8"

	LoggedForOffline  D8TxStatus = "9"  // записана для офлайн-обработки
	OfflineApplied    D8TxStatus = "10" // офлайн, применена к счёту
	OfflineNotApplied D8TxStatus = "11" // офлайн, применить к счёту не удалось

	GotRev           D8TxStatus = "12" // завершена, отклонена, получен реверс
	GotPrtRev        D8TxStatus = "13" // завершена, отклонена, получен частичный реверс
	InProgressGotRev D8TxStatus = "14"
	GotPartialRev    D8TxStatus = "15"

	FirstAttmptContinue D8TxStatus = "16"
	AdviceLogNotProceed D8TxStatus = "17"
)

type TxnType string

const (
	Sales   TxnType = "SALES"
	Cash    TxnType = "CASH"
	Deposit TxnType = "DEPOSIT"
	Balance TxnType = "BALANCE"
	C2C     TxnType = "TRANSF_C2C"
	H2C     TxnType = "TRANSF_H2C"
	C2A     TxnType = "TRANSF_C2A"
	A2C     TxnType = "TRANSF_A2C"
	Accver  TxnType = "ACCVER"
	Refund  TxnType = "REFUND"
	PreAuth TxnType = "PREAUTH"
)

const TJSCurrency = "972"
