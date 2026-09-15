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

var CardTypes = map[int]string{
	0: "1", // пластиковая;
	1: "2", //	TelebankID;
	2: "3", //	виртуальная
}

var CardStatuses = map[string]string{
	"": "unknown",

	"00": "1",  // Normal, active -> Open
	"01": "0",  // Card data prepared -> Not active
	"02": "0",  // Card data extracted -> Not active
	"03": "12", // Card prepared -> Not active
	"04": "0",  // Card production fail -> Not active
	"05": "5",  // VIP -> VIP
	"06": "6",  // Open Domestic
	"08": "8",  // Compromised
	"10": "4",  // PIN tries exceeded -> Restricted (ограничена)
	"11": "15", // Card expired -> Expired
	"12": "2",  // Card reported lost -> Lost
	"13": "3",  // Card reported stolen -> Stolen
	"14": "9",  // Customer closed -> Closed
	"15": "9",  // Bank cancelled -> Closed (закрыта банком)
	"16": "8",  // Card used fraudulent -> Compromised (заподозрена в мошенничестве)
	"17": "10", // Referral -> Referral (Работает c запросом к эмитенту)
	"20": "1",  //  ATM Operator card -> Open (активна, как ATM оператора)
}

// ReverseCardStatuses - обратный маппинг: TWO код -> список внешних кодов
var ReverseCardStatuses = map[string][]string{
	"0":  {"14"},       // Not active
	"1":  {"00"},       // Open
	"2":  {"12"},       // Lost
	"3":  {"13"},       // Stolen
	"4":  {"10"},       // Restricted: зачисление разрешено, расход запрещён
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
