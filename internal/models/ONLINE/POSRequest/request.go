package posrequestrq

import (
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"strconv"
)

type Body struct {
	SoapRq SoapRq `xml:"POSRequestRq" json:"POSRequestRq"`
}

func (b *Body) Call() (err error) {
	err = PosReq(b)
	return
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

	TranType   string   `xml:"TranType" json:"tran_type"`
	TranCode   TranCode `xml:"TranCode" json:"tran_code"`
	TranNumber string   `xml:"TranNumber" json:"tran_number"`

	TermName     string `xml:"TermName" json:"term_name"`
	TermInstName string `xml:"TermInstName" json:"term_inst_name"`

	PAN string `xml:"PAN" json:"pan"`

	FromAccount  string `xml:"FromAccount" json:"from_account"`
	FromAcctType string `xml:"FromAcctType" json:"from_acct_type"`
	ToAccount    string `xml:"ToAccount" json:"to_account"`
	ToAcctType   string `xml:"ToAcctType" json:"to_acct_type"`
	Amount       string `xml:"Amount" json:"amount"`

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

type TranCode string

const (
	Credit TranCode = "140"
	Debit  TranCode = "174"
)

func (req Request) GetTxnType() utils.TxnType {
	switch req.TranCode {
	case Credit:
		return utils.C2A
	case Debit:
		return utils.A2C
	default:
		return utils.Sales
	}
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

func (req Request) GetRecipientAcc() string {
	return req.ToAccount
}

func (req Request) GetTerminal() string {
	return req.TermName
}

func (req Request) GetAmount() float64 {
	amt, _ := strconv.ParseFloat(req.Amount, 64)
	return amt
}
