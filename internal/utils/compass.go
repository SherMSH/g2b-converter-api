package utils

import "strconv"

var AccountTypes = map[string]string{
	"":   "unknown",
	"00": "1",  // Checking (Расчётный / Текущий счёт)
	"11": "11", // Savings (Сберегательный / Накопительный счёт)
	"31": "31", // Credit (Кредитный счёт)
	"91": "91", // Bonus (Бонусный / Кешбэк-счёт)
}

var AccountStatuses = map[string]string{
	"": "unknown",

	"00": "1", // 1 – Open;
	"01": "0", // 0 – Inactive account;
	"02": "2", // 2 – Deposit only;
	"03": "3", // 3 – Open primary account;
	"04": "4", // 4 – Deposit only primary account;
	"09": "9", // 9 – Closed
	"10": "5", // 5 – Information only;
}

// ReverseAccountStatuses - обратный маппинг статуса счёта: код TWO -> код D8.
// Используется при установке статуса счёта через MDI.
var ReverseAccountStatuses = map[string]string{
	"1": "00", // Open
	"0": "01", // Inactive
	"2": "02", // Deposit only
	"3": "03", // Open primary account
	"4": "04", // Deposit only primary account
	"9": "09", // Closed
	"5": "10", // Information only
}

var CardTypes = map[int]string{
	0: "1", // пластиковая;
	1: "2", //	TelebankID;
	2: "3", //	виртуальная
}

// CardStatuses - перевод статуса карты D8 в статус TWO.
//
// Справочник статусов предоставлен командой процессинга: в спецификации их нет,
// поле statCode помечено как configurable value. В комментариях указан ответ
// авторизатора, который D8 отдаёт по каждому статусу - по нему видно, какие
// статусы процессинг пропускает, а какие отклоняет сам.
var CardStatuses = map[string]string{
	"": "unknown",

	"21": "0", // V-ACT P-PREP INA     0/00 - одобряет
	"22": "0", // V-ACT P-EXTR INA     0/00 - одобряет
	"23": "0", // V-ACT P-PROD OK INA  0/00 - одобряет
	"24": "0", // V-ACT P-PROD FAIL    0/00 - одобряет

	"00": "1",  // Normal, active        0/00 - одобряет
	"01": "0",  // Card data prepared    1/00 - отклоняет
	"02": "0",  // Card data extracted   1/00 - отклоняет
	"03": "12", // Card prepared         1/00 - отклоняет
	"04": "0",  // Card production fail  1/00 - отклоняет
	"05": "5",  // VIP                   0/03 - одобряет
	"06": "6",  // Open Domestic         0/05 - одобряет
	"08": "8",  // Compromised           1/00 - отклоняет общим кодом
	"10": "4",  // PIN tries exceeded    1/06 - отклоняет ВСЕ операции, включая зачисление
	"11": "15", // Card expired          1/01 - отклоняет
	"12": "2",  // Card reported lost    2/08 - отклоняет с изъятием
	"13": "3",  // Card reported stolen  2/09 - отклоняет с изъятием
	"14": "9",  // Customer closed       2/00 - отклоняет с изъятием
	"15": "9",  // Bank cancelled        2/00 - отклоняет с изъятием
	"16": "8",  // Card used fraudulent  2/02 - отклоняет с изъятием
	"17": "10", // Referral              0/01 - ОДОБРЯЕТ, хотя TWO требует отказ 71
	"20": "1",  // ATM Operator card     0/00 - одобряет
	"25": "4",  // Restricted            1/04 - отказ без изъятия, возвраты разрешены
}

// ReverseCardStatuses - обратный маппинг: TWO код -> список статусов D8.
// Используется при установке статуса, берётся первый код списка.
//
// Для статуса 4 (Restricted) в процессинге заведён отдельный statcode 25:
// отказ без изъятия карты с кодом 1/04 и разрешёнными возвратами. Прежний
// вариант 10 (PIN tries exceeded) оставлен запасным - он отклоняет все
// операции, включая зачисление.
var ReverseCardStatuses = map[string][]string{
	"0":  {"14"},       // Not active
	"1":  {"00"},       // Open
	"2":  {"12"},       // Lost
	"3":  {"13"},       // Stolen
	"4":  {"25", "10"}, // Restricted: зачисление разрешено, расход запрещён
	"5":  {"05"},       // VIP
	"6":  {"06"},       // Open Domestic
	"8":  {"08", "16"}, // Compromised: скомпрометирована либо использована мошеннически
	"9":  {"15"},       // Closed
	"10": {"17"},       // Referral (Необходим дополнительный запрос к эмитенту)
	"12": {"03"},       // Declared (Не издана)
	"15": {"11"},       // Expired
}

var Currencies = map[string]string{
	"":    "unknown",
	"972": "TJS",
	"978": "EUR",
	"840": "USD",
	"156": "CNY",
	"643": "RUB",
}

var TranCodes = map[int]string{
	0: "175", //Goods and services
	1: "175", //  01 Withdrawal
	// 2 Debit Adjustment
	// 9 Goods and services with cash disbursement
	// 10 Non-cash instrument
	// 11 Quasi cash
	// 17 Goods/sale with tip
	// 20 Refund
	21: "140", // 21 Deposits
	// 22 Credit Adjustment
	// 26 Cardholder funds transfer
	// 28 Cash Deposit (cash in)
	// 30 Available Funds Enquiry
	31: "117", // Balance Enquiry
	47: "135", // Money Transfer
	// 50 Bill payment
	// 58 Payment external bank
	// 90 PIN Change
	93: "116", // Customer Authentication
	// 94 PIN unblock (EMV only)
	// 95 Application unblock
}

// Currency переводит числовой код валюты в буквенный.
// Для пустого кода возвращает пустую строку, а не заглушку "unknown":
// внутренние литералы справочника не должны утекать партнёру.
func Currency(code string) string {
	if code == "" {
		return ""
	}
	return Currencies[code]
}

// TranCode переводит код операции D8 в код TWO.
// Коды, которых нет в TranCodes, отдаются как есть - терять их хуже,
// чем отдать партнёру незнакомое значение.
func TranCode(code int) string {
	if v, ok := TranCodes[code]; ok {
		return v
	}
	return strconv.Itoa(code)
}
