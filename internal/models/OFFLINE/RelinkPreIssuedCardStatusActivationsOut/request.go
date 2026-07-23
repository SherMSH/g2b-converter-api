package relinkpreissuedcardstatusactivationsout

import (
	models "converterapi/internal/models/OFFLINE"
	"converterapi/internal/utils"
	"encoding/xml"
)

// Root - корневой элемент XML
type Root struct {
	XMLName xml.Name         `xml:"ROOT"`
	Records []models.MRecord `xml:"RECORD"`
}

func (r Root) GetReqType() string {
	return string(utils.RelinkPreIssuedCardStatusActivationsOut)
}

func (r Root) GetRecords() []models.MRecord {
	return r.Records
}
func (r Root) GetRecordsCount() int {
	return len(r.Records)
}

func (r Root) Call() (respContent []byte, err error) {
	// mdiData, err := service.UpdatePreissiedCardG2b(r)
	// if err != nil {
	// 	return nil, err
	// }

	// for i := range mdiData.Details {
	// 	if mdiData.Details[i].C_ACTIONCODE != "0" {
	// 		break
	// 	}
	// 	pck := models.Pack{
	// 		CustomerId:   r.Records[i].PCode,
	// 		CustomerCode: r.Records[i].ExtID,
	// 		AccNum:       r.Records[i].Account,
	// 		CurrencyCode: "972",
	// 		LkeyAlias:    r.Records[i].ExternalID,
	// 		CardPan:      r.Records[i].PAN,
	// 	}
	// 	respContent = append(respContent, pck.GetData()...)
	// }
	return respContent, nil
}

// Record - запись с данными о карте
type Record struct {
	// Основные данные карты
	PAN        string `xml:"PAN"`        // Маскированный номер карты (9762****8427)
	MBR        string `xml:"MBR"`        // Member number (0)
	ExternalID string `xml:"EXTERNALID"` // Внешний ID (90263)

	// Статусы карты
	SignStat string `xml:"SIGNSTAT"` // Статус подписи (4)
	CRDStat  string `xml:"CRDSTAT"`  // Статус карты (1)

	// Данные на карте
	NameOnCard string `xml:"NAMEONCARD"` // Имя на карте (S.MAHMADIYAROVA)

	// CVV коды
	CVV  string `xml:"CVV"`  // CVV1 (пусто)
	CVV2 string `xml:"CVV2"` // CVV2 (пусто)
	IPVV string `xml:"IPVV"` // IPVV (пусто)

	// Дата окончания (пустая!)
	CancelDate string `xml:"CANCELDATE"` // Пустой тег

	// Валюты и риски
	CurrencyNo string `xml:"CURRENCYNO"` // Код валюты (пусто)
	RiskLevel  string `xml:"RISKLEVEL"`  // Уровень риска (пусто)

	// Дополнительные данные
	UserData   string `xml:"USERDATA"`   // Пользовательские данные (пусто)
	BRPart     string `xml:"BRPART"`     // Код отделения (27)
	LimitCMD   string `xml:"LIMITCMD"`   // Команда лимитов (пусто)
	BlockReiss string `xml:"BLOCKREISS"` // Блокировка перевыпуска (пусто)
	FinProfExt string `xml:"FINPROFEXT"` // Внешний финансовый профиль (пусто)
	ECStatus   string `xml:"ECSTATUS"`   // Статус EC (пусто)
	GroupCMD   string `xml:"GROUPCMD"`   // Команда группы (пусто)
	FinProfCMD string `xml:"FINPROFCMD"` // Команда финансового профиля (пусто)
	PinOffset  string `xml:"PINOFFSET"`  // PIN оффсет (пусто)
}
