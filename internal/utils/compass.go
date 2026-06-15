package utils

var AccountTypes = map[string]string{
	"":    "unknown",
	"00":  "1",  // Checking (Расчётный / Текущий счёт)
	"-01": "11", // Savings (Сберегательный / Накопительный счёт)
	"-02": "31", // Credit (Кредитный счёт)
	"-03": "91", // Bonus (Бонусный / Кешбэк-счёт)
}

var AccountStatuses = map[string]string{
	"":    "unknown",
	"-01": "0", // 0 – Inactive account;
	"01":  "1", // 1 – Open;
	"-02": "2", // 2 – Deposit only;
	"00":  "3", // 3 – Open primary account;
	"-04": "4", // 4 – Deposit only primary account;
	"-05": "5", // 5 – Information only;
	"-06": "9", // 9 – Closed
}

var CardTypes = map[int]string{
	0:  "1", // пластиковая;
	-1: "2", //	TelebankID;
	-2: "3", //	виртуальная
}

var Currencies = map[string]string{
	"":    "unknown",
	"972": "TJS",
	"978": "EUR",
	"840": "USD",
	"156": "CNY",
	"643": "RUB",
}
