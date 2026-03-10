package relinkpreissuedcardstatusactivationsout

import (
	"encoding/xml"
)

// Root - корневой элемент XML
type Root struct {
	XMLName xml.Name `xml:"ROOT"`
	Record  []Record `xml:"RECORD"`
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
