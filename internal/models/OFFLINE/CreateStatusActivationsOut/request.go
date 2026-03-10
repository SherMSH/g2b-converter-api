package createstatusactivationsout

import "encoding/xml"

// Root - корневой элемент XML
type Root struct {
	XMLName xml.Name `xml:"ROOT"`
	Record  []Record `xml:"RECORD"`
}

// Record - запись с данными о карте
type Record struct {
	// Основные данные карты
	PAN        string `xml:"PAN"`        // Маскированный номер карты
	MBR        string `xml:"MBR"`        // Member number (0)
	ExternalID string `xml:"EXTERNALID"` // Внешний ID (161995)

	// Статусы карты
	SignStat string `xml:"SIGNSTAT"` // Статус подписи (4)
	CRDStat  string `xml:"CRDSTAT"`  // Статус карты (1)

	// Данные на карте
	NameOnCard string `xml:"NAMEONCARD"` // Имя на карте

	// CVV коды
	CVV  string `xml:"CVV"`  // CVV1
	CVV2 string `xml:"CVV2"` // CVV2
	IPVV string `xml:"IPVV"` // IPVV (Integrated PIN Verification Value)

	// Даты
	CancelDate string `xml:"CANCELDATE"` // Дата окончания действия (29022028)

	// Валюты и риски
	CurrencyNo string `xml:"CURRENCYNO"` // Код валюты
	RiskLevel  string `xml:"RISKLEVEL"`  // Уровень риска

	// Дополнительные данные
	UserData   string `xml:"USERDATA"`   // Пользовательские данные
	BRPart     string `xml:"BRPART"`     // Код отделения (53)
	LimitCMD   string `xml:"LIMITCMD"`   // Команда лимитов
	BlockReiss string `xml:"BLOCKREISS"` // Блокировка перевыпуска
	FinProfExt string `xml:"FINPROFEXT"` // Внешний финансовый профиль
	ECStatus   string `xml:"ECSTATUS"`   // Статус EC (Electronic Commerce)
	GroupCMD   string `xml:"GROUPCMD"`   // Команда группы
	FinProfCMD string `xml:"FINPROFCMD"` // Команда финансового профиля
	PinOffset  string `xml:"PINOFFSET"`  // PIN оффсет
}
