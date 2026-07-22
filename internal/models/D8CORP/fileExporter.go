package d8corp

type MdiData struct {
	Header  RespHeader    `json:"mdiResponseHeader"`
	Details []RespDetails `json:"mdiResponseDetails"`
}

type RespHeader struct {
	CActionCode  string `json:"C_ACTIONCODE"`
	CFileName    string `json:"C_FILE_NAME"`
	CRspCode     string `json:"C_RSPCODE"`
	IssRecAction string `json:"ISS_RECACTION"`
	IssRecType   string `json:"ISS_RECTYPE"`
	IssTimestamp string `json:"ISS_TIMESTAMP"`
	IRejMsg      string `json:"I_REJMSG"`
}

type RespDetails struct {
	C_ACTIONCODE  string `json:"C_ACTIONCODE"`
	C_RSPCODE     string `json:"C_RSPCODE"`
	ISS_RECACTION string `json:"ISS_RECACTION"`
	ISS_RECNUM    int    `json:"ISS_RECNUM"`
	ISS_RECTYPE   string `json:"ISS_RECTYPE"`
	I_REJMSG      string `json:"I_REJMSG"`

	ISS_CARDA_ID  int    `json:"ISS_CARDA_ID"`
	ISS_CARD_ID   int    `json:"ISS_CARD_ID"`
	KL_LKEY_ALIAS string `json:"KL_LKEY_ALIAS"`
	KL_LKEY_CLR   string `json:"KL_LKEY_CLR"`
	KL_LKEY_SEQNO int    `json:"KL_LKEY_SEQNO"`
}
