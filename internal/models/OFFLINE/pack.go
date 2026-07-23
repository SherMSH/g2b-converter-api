package models

import "fmt"

type Pack struct {
	CustomerId   string
	CustomerCode string
	AccNum       string
	CurrencyCode string
	LkeyAlias    string
	CardPan      string
}

func (p *Pack) GetData() []byte {
	return []byte(fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s\n", p.CustomerId, p.CustomerCode, p.AccNum, p.AccNum, p.CurrencyCode, p.LkeyAlias, p.CardPan))
}
