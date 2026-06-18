package utils

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
	"03": "0",  // Card prepared -> Not active
	"04": "0",  // Card production fail -> Not active
	"10": "4",  // PIN tries exceeded -> Restricted
	"11": "15", // Card expired -> Expired
	"12": "2",  // Card reported lost -> Lost
	"13": "3",  // Card reported stolen -> Stolen
	"14": "9",  // Customer closed -> Closed
	"15": "9",  // Bank cancelled -> Closed (закрыта банком)
	"16": "4",  // Card used fraudulent -> Restricted (скомпрометирована/мошенничество)
	"20": "1",  //  ATM Operator card -> Open (активна, как ATM оператора)
}

// ReverseCardStatuses - обратный маппинг: TWO код -> список внешних кодов
var ReverseCardStatuses = map[string][]string{
	"0":  {"01", "02", "03", "04"}, // Not active
	"1":  {"00", "20"},             // Open
	"2":  {"12"},                   // Lost
	"3":  {"13"},                   // Stolen
	"4":  {"10", "16"},             // Restricted
	"5":  {},                       // VIP (нет прямого маппинга)
	"6":  {},                       // Open Domestic (нет прямого маппинга)
	"8":  {},                       // Compromised (нет прямого маппинга)
	"9":  {"14", "15"},             // Closed
	"10": {},                       // Referral (нет прямого маппинга)
	"12": {},                       // Declared (нет прямого маппинга)
	"15": {"11"},                   // Expired
}

var Currencies = map[string]string{
	"":    "unknown",
	"972": "TJS",
	"978": "EUR",
	"840": "USD",
	"156": "CNY",
	"643": "RUB",
}
