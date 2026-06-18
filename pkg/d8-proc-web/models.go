package d8procweb

type CardData struct {
	CardaID       int    `json:"carda_id"`
	CardID        int    `json:"card_id"`
	CardaRecver   int    `json:"carda_recver"`
	CardRecver    int    `json:"card_recver"`
	LkeyID        int    `json:"lkey_id"`
	LkeyDisplay   string `json:"lkey_display"`
	Alias         string `json:"alias"`
	CompanyID     int    `json:"company_id"`
	Name          string `json:"name"`
	Custcode      string `json:"custcode"`
	Embossname    string `json:"embossname"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	DelyveryPoint string `json:"delyvery_point"`
	Currcode      string `json:"currcode"`
	Expdate       string `json:"expdate"`
	CompanyRegnum string `json:"company_regnum"`
	CompanyName   string `json:"company_name"`
	Bincode       string `json:"bincode"`
	BinID         int    `json:"bin_id"`

	PAN         string `json:"-"`
	ProductType int    `json:"-"`
	StatCode    string `json:"-"`
}

type AccountData struct {
	ID          int     `json:"account_id"`
	Recver      int     `json:"account_recver,omitempty"`
	CompanyID   int     `json:"company_id,omitempty"`
	CustomerID  int     `json:"customer_id,omitempty"`
	Typecode    string  `json:"typecode,omitempty"`
	Currcode    string  `json:"currcode,omitempty"`
	Accnum      string  `json:"accnum,omitempty"`
	AvlBal      float64 `json:"avlbal,omitempty"`
	BlkAmt      float64 `json:"blkamt,omitempty"`
	AvlbalUnset float64 `json:"avlbal_unset,omitempty"`
	BlkamtUnset float64 `json:"blkamt_usnset,omitempty"`
	OpenDate    string  `json:"opendate,omitempty"`
	LastUsage   string  `json:"lastusage,omitempty"`
	Statcode    string  `json:"statcode,omitempty"`
	Crlimit     float64 `json:"crlimit,omitempty"`
	Balincr     float64 `json:"balincr,omitempty"`
	BalincrExp  string  `json:"balincrexp,omitempty"`

	Name          string `json:"name,omitempty"`
	Custcode      string `json:"custcode,omitempty"`
	CompanyRegnum string `json:"company_regnum,omitempty"`
	CompanyName   string `json:"company_name,omitempty"`
}

type Filter struct {
	Column string   `json:"column"`
	Values []string `json:"values"`
}

type RequestBody struct {
	Data     interface{} `json:"data,omitempty"`
	Ordering []string    `json:"ordering,omitempty"`
	Filters  []Filter    `json:"filters,omitempty"`
}
