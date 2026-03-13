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
	RemoveCMSAbonentRq       RqBodyType = "RemoveCMSAbonentRq"
	RemovePersonCMSAbonentRq RqBodyType = "RemovePersonCMSAbonentRq"
	ResetBadPINTriesRq       RqBodyType = "ResetBadPINTriesRq"
	SetCardStatusRq          RqBodyType = "SetCardStatusRq"
	UpdatePersonRq           RqBodyType = "UpdatePersonRq"
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
	RemoveCMSAbonentRq,
	RemovePersonCMSAbonentRq,
	ResetBadPINTriesRq,
	SetCardStatusRq,
	UpdatePersonRq,
}

type OfflineReqType string

const (
	CreateCardsOut                          OfflineReqType = "CreateCardsOut.xml"
	CreateCustomerAndAccount                OfflineReqType = "CreateCustomerAndAccount.xml"
	CreateOrganizations                     OfflineReqType = "CreateOrganizations.xml"
	CreatePreIssuedCards                    OfflineReqType = "CreatePreIssuedCards.xml"
	CreateStatusActivationsOut              OfflineReqType = "CreateStatusActivationsOut.xml"
	ReissueCardsOut                         OfflineReqType = "ReissueCardsOut.xml"
	RelinkPreIssuedCardsOut                 OfflineReqType = "RelinkPreIssuedCardsOut.xml"
	RelinkPreIssuedCardStatusActivationsOut OfflineReqType = "RelinkPreIssuedCardStatusActivationsOut.xml"
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
