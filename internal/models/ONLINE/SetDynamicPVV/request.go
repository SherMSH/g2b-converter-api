package setdynamicpvv

import "reflect"

type Body struct {
	SoapRq SoapRq `xml:"SetDynamicPVV_PINOffsetRq" json:"SetDynamicPVV_PINOffsetRq"`
}

func (sb Body) GetBodyType() reflect.Type {
	return reflect.TypeOf(sb)
}

func (sb *Body) Call() (*Envelope, error) {
	return Svc(sb)
}

// SoapRq соответствует элементу fimi:SetDynamicPVV_PINOffsetRq
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

	PAN     string `xml:"PAN" json:"pan"`
	MBR     string `xml:"MBR" json:"mbr"`
	CardUID string `xml:"CardUID" json:"card_uid"`

	// PINBlock - значение PIN. Процессинг D8 принимает PIN-блок только под
	// одноразовым ключом, который собираем мы сами, поэтому здесь ожидается
	// открытый PIN, а не блок под TPK - см. комментарий в Svc.
	PINBlock string `xml:"PINBlock" json:"pin_block"`

	// KeyId - идентификатор рабочего ключа. В схеме D8 рабочих ключей нет,
	// поле принимается для совместимости и не используется.
	KeyId string `xml:"KeyId" json:"key_id"`

	// PVKI - индекс ключа проверки PIN. D8 задаёт его на уровне карты, через
	// setPIN не меняется; поле принимается, но не применяется.
	PVKI string `xml:"PVKI" json:"pvki"`

	// ExpDate в формате YYMM. Если не прислан - берём из карты
	ExpDate string `xml:"ExpDate" json:"exp_date"`

	// SingleOperation - PIN на одну операцию. D8 такого режима не имеет:
	// setPIN меняет PIN насовсем.
	SingleOperation string `xml:"SingleOperation" json:"single_operation"`

	ChangeReason string `xml:"ChangeReason" json:"change_reason"`
}
