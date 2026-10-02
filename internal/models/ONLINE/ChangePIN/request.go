package changepin

type Body struct {
	SoapRq SoapRq `xml:"ChangePIN" json:"ChangePIN"`
}

// RqBody - то же тело под тегом ChangePINRq. Партнёр пишет операцию и так, и
// так, поэтому принимаем оба написания.
type RqBody struct {
	SoapRq SoapRq `xml:"ChangePINRq" json:"ChangePINRq"`
}

func (sb *RqBody) Call() (*Envelope, error) {
	return Svc(&Body{SoapRq: sb.SoapRq})
}

func (sb *Body) Call() (*Envelope, error) {
	rsp, err := Svc(sb)
	return rsp, err
}

// SoapRq соответствует элементу ChangePIN
type SoapRq struct {
	Req Request `xml:"Request" json:"Request"`
}

// Request соответствует элементу fimi:Request
type Request struct {
	Ver      string `xml:"Ver,attr" json:"ver"`
	Product  string `xml:"Product,attr" json:"product"`
	Echo     string `xml:"Echo,attr" json:"echo"`
	Session  string `xml:"Session,attr" json:"session"`
	Clerk    string `xml:"Clerk,attr" json:"clerk"`
	Password string `xml:"Password,attr" json:"password"`

	PAN string `xml:"PAN" json:"pan"`
	MBR string `xml:"MBR" json:"mbr"`
}
