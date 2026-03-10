package reissuecardsout

import "encoding/xml"

// Root - корневой элемент XML
type Root struct {
	XMLName xml.Name `xml:"ROOT"`
	Record  []Record `xml:"RECORD"`
}

// Record - запись с данными о карте
type Record struct {
	// Основные данные карты
	PAN        string `xml:"PAN"`        // Маскированный номер карты (9762****7014)
	MBR        string `xml:"MBR"`        // Member number (0)
	ExternalID string `xml:"EXTERNALID"` // Внешний ID (109730_142982)

	// Продукт и причина
	CardProd string `xml:"CARDPROD"` // Продукт карты (4)
	Reason   string `xml:"REASON"`   // Причина операции (2)

	// Флаги и приоритеты
	MakePrior   string `xml:"MAKEPRIOR"`   // Приоритет изготовления (0)
	CMChangePAN string `xml:"CMCHANGEPAN"` // Флаг смены PAN (1)

	// Данные на карте
	NameOnCard string `xml:"NAMEONCARD"` // Имя на карте (M.RAKHMATOV)

	// Финансовые профили
	FinProfile string `xml:"FINPROFILE"` // Финансовый профиль
	FinProfExt string `xml:"FINPROFEXT"` // Внешний финансовый профиль

	// Даты
	CancelDate string `xml:"CANCELDATE"` // Дата окончания действия (30092027)

	// Переизготовление
	RemakePAN string `xml:"REMAKEPAN"` // PAN при переизготовлении
	RemakeMBR string `xml:"REMAKEMBR"` // MBR при переизготовлении

	// Отделение
	BRPart string `xml:"BRPART"` // Код отделения (22)
}
