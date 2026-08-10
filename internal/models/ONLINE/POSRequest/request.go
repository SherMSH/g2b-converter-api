package posrequestrq

import (
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
)

type Body struct {
	SoapRq SoapRq `xml:"POSRequestRq" json:"POSRequestRq"`
}

func (b *Body) Call() (resp *Envelope, err error) {
	return PosReq(b)
}

func (b Body) GetReqType() interface{} {
	return utils.POSRequestRq
}

// SoapRq соответствует элементу GetAcctInfoRq
type SoapRq struct {
	Req          Request `xml:"Request" json:"request"`
	ApprovalCode string  `xml:"-" json:"-"`
}

// Request соответствует элементу Request
type Request struct {
	Ver      string `xml:"Ver,attr" json:"ver"`
	Product  string `xml:"Product,attr" json:"product"`
	Echo     string `xml:"Echo,attr" json:"echo"`
	Clerk    string `xml:"Clerk,attr" json:"clerk"`
	Password string `xml:"Password,attr" json:"password"`

	TranType   TranType `xml:"TranType" json:"tran_type"`
	TranCode   TranCode `xml:"TranCode" json:"tran_code"`
	TranNumber string   `xml:"TranNumber" json:"tran_number"`

	TermName     string `xml:"TermName" json:"term_name"`
	TermInstName string `xml:"TermInstName" json:"term_inst_name"`

	PAN string `xml:"PAN" json:"pan"`

	// Карта получателя в переводах: FIMI передаёт её вторым набором полей
	PAN2 string `xml:"PAN2" json:"pan2"`
	MBR2 string `xml:"MBR2" json:"mbr2"`

	// Назначение перевода (Appendix D спецификации D8): если партнёр не прислал,
	// подставляется значение по типу операции
	BusinessApplicationIdent string `xml:"BusinessApplicationIdent" json:"business_application_ident"`

	FromAccount  string  `xml:"FromAccount" json:"from_account"`
	FromAcctType string  `xml:"FromAcctType" json:"from_acct_type"`
	ToAccount    string  `xml:"ToAccount" json:"to_account"`
	ToAcctType   string  `xml:"ToAcctType" json:"to_acct_type"`
	Amount       float64 `xml:"Amount" json:"amount"`

	CVV       string `xml:"CVV" json:"cvv"`
	CVV2      string `xml:"CVV2" json:"cvv2"`
	EntryMode string `xml:"EntryMode" json:"entry_mode"`
	Condition string `xml:"Condition" json:"condition"`
	Currency  string `xml:"Currency" json:"currency"`
	Track2    string `xml:"Track2" json:"track2"`

	ThisTranId string `xml:"ThisTranId" json:"this_tran_id"`
	ECTranId   string `xml:"ECTranId" json:"ec_tran_id"`
	OrigTranId string `xml:"OrigTranId" json:"orig_tran_id"`

	MBR      string `xml:"MBR" json:"mbr"`
	RespCode string `xml:"RespCode" json:"resp_code"`
}

type TranType int

// Тип транзакции:
// 100=AuthRequest (авторизационный
// запрос, происходит онлайн-авторизация);
// 120=AuthAdvice (авторизационный запрос,
// безусловный прием и одобрение
// извещения о транзакции);
// 200=Request (финансовый запрос,
// происходит онлайн-авторизация);
// 220=Advice (финансовый запрос,
// безусловный прием и одобрение
// извещения о транзакции);
// 999=Admin(админстративный запрос,
// используется только для транзакции
// CloseDay (162))

type TranCode int

// Код транзакции. Возможные значения:
// - Purchase (110),
// - Pre-purchase(111),
// - Pre-purchase Complete (112),
// - Mail/Phone Order(113),
// - Merchandise Return(114),
// - Cash Advance(115),
// - Card Verification(116),
// - POS Balance Inquiry(117),
// - Purchase With Cashback(118),
// - Purchase Adjustment(121),
// - Merchandise Return Adjustment(122),
// - Cash Advance Adjustment(123),
// - Pre-purchase Increment(124),
// - Purchase Cancellation(125),
// - POS Message to Financial Institution(126),
// - Quasi-Cash(130),
// - POS P2P Debit(132),
// - POS P2P Credit(133),
// - POS P2P Calc Fee(134),
// - POS P2P Pass(135),
// - POS Installment Details(136),
// - POS PIN Change(139),
// - POS Deposit(140),
// - POS Deposit Adjustment(141),
// - POS Transfer Pass(149),
// - Close Batch(160),
// - Close Day(162),
// - POS Prepaid Pass(171),

const (
	Credit   TranCode = 140
	P2P      TranCode = 135
	Transfer TranCode = 149
	Debit    TranCode = 175
)

// GetTxnType определяет тип операции D8 по коду операции FIMI.
//
// Для переводов конкретный тип зависит от того, чем заданы стороны:
// карта получателя (PAN2) или счёт (ToAccount/FromAccount).
//
// Неизвестный код возвращает пустое значение - авторизация такую операцию
// отклонит. Раньше любой нераспознанный код молча превращался в SALES, то есть
// в списание.
func (req Request) GetTxnType() utils.TxnType {
	switch req.TranCode {
	case Credit:
		return utils.Deposit
	case Debit:
		return utils.Sales
	case P2P, Transfer:
		switch {
		case req.PAN2 != "" && req.FromAccount != "":
			return utils.A2C
		case req.ToAccount != "":
			return utils.C2A
		case req.PAN2 != "":
			return utils.C2C
		default:
			return ""
		}
	default:
		return ""
	}
}

// Реквизиты перевода: реализация models.TrnTransferIface

func (req Request) GetRecipientPan() string {
	return req.PAN2
}

func (req Request) GetRecipientAccount() string {
	return req.ToAccount
}

func (req Request) GetSenderAccount() string {
	return req.FromAccount
}

func (req Request) GetDestinationAccountType() string {
	return req.ToAcctType
}

func (req Request) GetBusinessAppId() string {
	return req.BusinessApplicationIdent
}

func (req Request) GetPan() string {
	return req.PAN
}

func (req Request) GetMBR() string {
	return req.MBR
}

func (req Request) GetExpDate() string {
	logger.Infof("req: %+v", req)
	return utils.GetExpFromTrack(req.Track2)
}

func (req Request) GetAccount() string {
	return req.FromAccount
}

func (req Request) GetCurrency() string {
	if len(req.Currency) == 0 {
		return utils.TJSCurrency
	}
	return req.Currency
}

func (req Request) GetTerminal() string {
	return req.TermName
}

func (req Request) GetAmount() float64 {
	return req.Amount
}

func (req Request) GetAcceptorID() string {
	terminals, err := service.GetTerminalByIdG2b(req.TermName)
	if err != nil || len(terminals) == 0 {
		return "000001"
	}
	return terminals[0].MerchantIDStr
}
