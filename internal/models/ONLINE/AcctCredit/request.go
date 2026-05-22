package acctcredit

import "reflect"

type Body struct {
	SoapRq SoapRq `xml:"AcctCreditRq" json:"AcctCreditRq"`
}

func (sb Body) GetBodyType() reflect.Type {
	return reflect.TypeOf(sb)
}

func (sb *Body) Call() (*Envelope, error) {
	return Svc(sb)
}

// SoapRq соответствует элементу GetAcctInfoRq
type SoapRq struct {
	Req Request `xml:"Request" json:"request"`
}

// Request соответствует элементу Request
type Request struct {
	Ver      string `xml:"Ver,attr" json:"ver"`
	Product  string `xml:"Product,attr" json:"product"`
	Clerk    string `xml:"Clerk,attr" json:"clerk"`
	Password string `xml:"Password,attr" json:"password"`

	Account          string `xml:"Account" json:"account"`
	InstName         string `xml:"InstName" json:"inst_name"`
	PersonId         string `xml:"PersonId" json:"person_id"`
	AccountUID       string `xml:"AccountUID" json:"account_uid"`
	AggregateId      string `xml:"AggregateId" json:"aggregate_id"`
	PAN              string `xml:"PAN" json:"pan"`
	MBR              string `xml:"MBR" json:"mbr"`
	CardUID          string `xml:"CardUID" json:"card_uid"`
	Amount           string `xml:"Amount" json:"amount"`
	CorrespAcct      string `xml:"CorrespAcct" json:"corresp_acct"`
	CorrespAcctUID   string `xml:"CorrespAcctUID" json:"corresp_acct_uid"`
	TranNumber       string `xml:"TranNumber" json:"tran_number"`
	OrigAmount       string `xml:"OrigAmount" json:"orig_amount"`
	OrigCurrency     string `xml:"OrigCurrency" json:"orig_currency"`
	IgnoreImpact     string `xml:"IgnoreImpact" json:"ignore_impact"`
	NeedNotify       string `xml:"NeedNotify" json:"need_notify"`
	PrevTranId       string `xml:"PrevTranId" json:"prev_tran_id"`
	RefreshSeqNo     string `xml:"RefreshSeqNo" json:"refresh_seq_no"`
	ChangeReason     string `xml:"ChangeReason" json:"change_reason"`
	UseBonusDebt     string `xml:"UseBonusDebt" json:"use_bonus_debt"`
	BonusProgramName string `xml:"BonusProgramName" json:"bonus_program_name"`
	PrizeID          string `xml:"PrizeID" json:"prize_id"`
	PrizeQuantity    string `xml:"PrizeQuantity" json:"prize_quantity"`
}
