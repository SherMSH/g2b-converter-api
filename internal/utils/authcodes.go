package utils

// Коды ответа авторизатора TWO.
//
// Партнёр ожидает коды из справочника TWO (документ «Коды ошибок online»), а
// процессинг D8 отвечает своей парой code/rspcode (Appendix A спецификации
// D8 G2B). Здесь собран перевод одного в другое.
const (
	TwoApproved        = "1"  // Approved - OK
	TwoApprovedPartial = "2"  // Одобрено на частичную сумму
	TwoExternalDecline = "68" // Отказ внешнего хоста: отказ, которого нет в таблице
	TwoSystemError     = "54" // Системная ошибка
	// TwoInsufficientFunds - нехватка средств: причину отказа по нему партнёр
	// ждёт с разбором по слагаемым остатка
	TwoInsufficientFunds = "59"
	TwoNone              = "0" // None - ответа авторизатора нет
)

// d8AuthRespCodes - перевод пары code/rspcode процессинга в код ответа TWO.
// Ключ - "code/rspcode".
//
// Слева коды из Appendix A спецификации D8, справа - из справочника
// «Основные коды ответа авторизатора (RespCode)» партнёра.
var d8AuthRespCodes = map[string]string{
	// Одобрения
	"0/00": TwoApproved,
	"0/01": TwoApproved,        // Honour with ID
	"0/02": TwoApprovedPartial, // Approved for a partial amount
	"0/80": TwoApproved,        // ICC offline approved
	"0/81": TwoApproved,        // ICC unable to go online, approved
	"0/82": TwoApproved,        // ICC approved after card-initiated referral
	"0/85": TwoApproved,        // Not declined, zero amount

	// Отказы
	"1/00": "50", // Do not honour -> Unauthorized usage
	"1/01": "51", // Card expired
	"1/02": "75", // Suspected fraud -> External decline special condition
	"1/03": "71", // Acceptor contact acquirer -> Contact card issuer
	"1/04": "58", // Restricted card
	"1/05": "71", // Acceptor contact security
	"1/06": "62", // PIN tries exceeded -> PIN tries limit was reached
	"1/07": "71", // Acceptor contact issuer
	"1/08": "71", // Refer to issuer special conditions
	"1/09": "98", // Invalid merchant
	"1/10": "67", // Invalid amount
	"1/11": "52", // Invalid card
	"1/12": "53", // PIN required -> Invalid PIN
	"1/13": "55", // Unacceptable fee -> Ineligible transaction
	"1/14": "56", // No account of specified type -> Ineligible account
	"1/15": "57", // Function not supported -> Transaction not supported
	"1/16": "59", // Insufficient funds
	"1/17": "53", // Bad PIN -> Invalid PIN
	"1/18": "52", // No card record -> Invalid card
	"1/19": "57", // Transaction not allowed to cardholder
	"1/20": "69", // Transaction not allowed to terminal -> No sharing
	"1/21": "61", // Exceeds amount limit -> Withdrawal limit would be exceeded
	"1/22": "50", // Security violation
	"1/23": "60", // Exceeds frequency limit -> Uses limit exceeded
	"1/24": "50", // Violation of law
	"1/25": "50", // Card not effective -> Unauthorized usage
	"1/26": "53", // Invalid PIN block
	"1/27": "53", // PIN length error
	"1/28": TwoSystemError,
	"1/29": "75", // Suspect counterfeit
	"1/65": "55", // AML -> Ineligible transaction
	"1/70": "71", // Cardholder contact issuer
	"1/71": "53", // PIN not changed
	"1/75": "61", // Exceeds limit: PIN required
	"1/76": "6",  // Exceeds limit: additional authentication -> Strong customer auth required
	"1/77": "67", // Invalid currency code -> Invalid amount
	"1/78": "50", // ICC decline after referral
	"1/79": "25", // New PIN unsafe -> Weak PIN
	"1/80": "67", // Amount not divisible
	"1/81": "67", // Amount too big for ATM
	"1/82": "95", // ATM out of cash -> Insufficient cash
	"1/83": "56", // No chequing account -> Ineligible account
	"1/84": "53", // New PINs entered differ
	"1/85": "81", // CVV validation error -> Bad CVV2
	"1/86": "82", // pre-auth time limit exceeded -> Invalid transaction
	"1/87": "50", // Card has been destroyed
	"1/88": "59", // Restricted card. No funds -> Insufficient funds
	"1/89": "55", // Foreign exchange not available
	"1/90": "60", // Transaction count limit exceeded
	"1/91": "72", // System overload -> Destination not available
	"1/92": "66", // Statement data not available
	"1/93": "83", // Pin tries exceeded (лимит уже был достигнут ранее)
	"1/94": "65", // No balance data available
	"1/95": "50", // ICC offline declined
	"1/96": "50", // ICC unable to go online, declined
	"1/97": "85", // ICC card authentication failed -> Bad ARQC

	// Изъятие карты
	"2/00": "50", // Decline, pickup card
	"2/01": "57", // Not supported by receiver
	"2/02": "75", // Suspected fraud
	"2/03": "71", // Acceptor contact acquirer
	"2/04": "58", // Restricted card
	"2/05": "71", // Acceptor contact security
	"2/06": "62", // PIN tries exceeded
	"2/07": "71", // Special conditions
	"2/08": "40", // Lost card
	"2/09": "41", // Stolen card
	"2/10": "75", // Suspect counterfeit
	"2/85": "81", // CVV validation error

	// Ошибки обработки. Не все из них системные: часть несёт обычную причину
	// отказа, и сваливать их в 54 было бы потерей смысла.
	"9/02": "82", // Invalid transaction
	"9/03": "82", // Re-enter transaction
	"9/04": "74", // Format error
	"9/05": "69", // Acquirer not supported -> No sharing
	"9/06": "23", // Cutover in progress
	"9/07": "72", // Issuer/switch inoperative -> Destination not available
	"9/08": "73", // Destination not found -> Routing error
	"9/09": TwoSystemError,
	"9/10": "72", // Issuer signed off
	"9/11": "72", // Issuer timed out
	"9/12": "72", // Issuer unavailable
	"9/13": "82", // Duplicate transaction
	"9/14": "15", // Unable to trace original transaction -> Original transaction not found
	"9/15": TwoSystemError,
	"9/16": TwoSystemError, // MAC incorrect
	"9/17": TwoSystemError, // MAC key sync error
	"9/18": TwoSystemError, // No comms keys available
	"9/19": TwoSystemError, // Encryption key sync error
	"9/20": "50",           // Security error, try again
	"9/21": "50",           // Security error, no action
	"9/22": TwoSystemError, // Message out of sequence
	"9/23": "4",            // Request in progress -> Postponed
	"9/40": "74",           // Invalid transaction date
	"9/50": "50",           // Disagreement
	"9/51": "51",           // Card expired
	"9/52": "75",           // Suspected fraud
	"9/53": "58",           // Restricted card
	"9/54": "98",           // Invalid merchant
	"9/55": "67",           // Invalid amount
	"9/56": "52",           // Invalid card number
	"9/57": "55",           // Unacceptable fee
	"9/58": "59",           // Insufficient funds
	"9/59": "53",           // Bad PIN
	"9/60": "57",           // Trans not allowed to cardholder
	"9/61": "69",           // Trans not allowed to terminal
	"9/62": "61",           // Exceeds amount limit
	"9/63": "50",           // Security violation
	"9/64": "60",           // Exceeds frequency limit
	"9/65": "50",           // Violation of law
	"9/66": TwoSystemError, // Reconciliation error
	"9/67": TwoSystemError, // MAC incorrect
	"9/68": TwoSystemError, // MAC key sync error
	"9/69": TwoSystemError, // No comms key available
	"9/70": TwoSystemError, // Encryption key sync error
	"9/71": "50",           // Security error, try again
	"9/72": "50",           // Security error, no action
	"9/73": TwoSystemError, // Message out of sequence

	// Служебные коды уровня сервиса
	"D/00": "52", // Card not found -> Invalid card
	"D/01": "50", // Card is in invalid status
	"D/02": "82", // Status change transition not allowed
	"D/03": "56", // Card not linked to Account -> Ineligible account
	"D/04": "15", // Transaction not found -> Original transaction not found
	"D/05": "52", // Lkey not found
	"D/06": "15", // Invalid ecTxRefno
	"D/07": "56", // Invalid account
	"C/42": "74", // Clear PAN not allowed in request -> Format error
	"C/43": "74", // CVV2 not allowed in request
	"C/44": TwoSystemError,
}

// AuthRespCode переводит ответ процессинга в код ответа авторизатора TWO.
//
// Неизвестный отказ отдаётся как 68 (External decline): партнёру важно понимать,
// что операция отклонена, а подменять причину выдуманной хуже, чем признать её
// внешней. Неизвестный код одобрения считается обычным одобрением.
func AuthRespCode(code, rspcode string) string {
	if v, ok := d8AuthRespCodes[code+"/"+rspcode]; ok {
		return v
	}
	switch code {
	case "0":
		return TwoApproved
	case "":
		return TwoNone
	default:
		// Сюда попадают неизвестные коды всех семейств, включая 9 (ошибки
		// обработки): известные из них разобраны в таблице выше.
		return TwoExternalDecline
	}
}

// IsApproved - признак одобрения операции процессингом.
func IsApproved(code string) bool {
	return code == "0"
}

// RetainCard - признак изъятия карты: код 2 в D8 означает pickup.
func RetainCard(code string) string {
	if code == "2" {
		return "1"
	}
	return "0"
}

// CardStatusRespCodes - код ответа авторизатора по статусу карты TWO.
// Из таблицы «Статусы карты в TWO и коды ответов авторизатора».
// Статусы 1 (Open), 5 (VIP) и 6 (Open Domestic) операции не запрещают.
var CardStatusRespCodes = map[string]string{
	"0":  "50", // Not active
	"2":  "40", // Lost
	"3":  "41", // Stolen
	"8":  "75", // Compromised
	"9":  "50", // Closed
	"10": "71", // Referral - нужен запрос к эмитенту
	"12": "50", // Declared - не издана
	"15": "51", // Expired
}

// CardStatusNames - названия статусов карты в TWO.
// Нужны для текста причины отказа: партнёр ждёт формулировку вида
// "Response for card status 'Lost' ...".
var CardStatusNames = map[string]string{
	"0":  "Not active",
	"1":  "Open",
	"2":  "Lost",
	"3":  "Stolen",
	"4":  "Restricted",
	"5":  "VIP",
	"6":  "Open Domestic",
	"8":  "Compromised",
	"9":  "Closed",
	"10": "Referral",
	"12": "Declared",
	"15": "Expired",
}

// AccountStatusNames - названия статусов счёта в TWO.
var AccountStatusNames = map[string]string{
	"0": "Inactive",
	"1": "Open",
	"2": "Deposit only",
	"3": "Open primary account",
	"4": "Deposit only primary account",
	"5": "Information only",
	"9": "Closed",
}

// CardStatusRespCode возвращает код отказа по статусу карты.
//
// Статус 4 (Restricted) разрешает зачисление и запрещает расход, поэтому
// зависит от направления операции - как и статусы счёта «только приход».
// Пустая строка означает, что статус карты операцию не запрещает.
func CardStatusRespCode(status string, isDebit bool) string {
	if code, ok := CardStatusRespCodes[status]; ok {
		return code
	}
	if isDebit && status == "4" {
		return "58" // Restricted Card
	}
	return ""
}

// AccountStatusRespCodes - код ответа авторизатора по статусу счёта TWO.
// Из таблицы «Статусы счёта карты в TWO и коды ответов авторизатора».
// Статусы 2 и 4 (Deposit only) запрещают только расход, поэтому их проверяет
// AccountStatusRespCode.
var AccountStatusRespCodes = map[string]string{
	"0": "56", // Inactive
	"5": "55", // Information only
	"9": "56", // Closed
}

// AccountStatusRespCode возвращает код отказа по статусу счёта, учитывая
// направление операции: «только приход» отклоняет лишь расходные операции.
// Пустая строка означает, что статус счёта операцию не запрещает.
func AccountStatusRespCode(status string, isDebit bool) string {
	if code, ok := AccountStatusRespCodes[status]; ok {
		return code
	}
	if isDebit && (status == "2" || status == "4") {
		return "55" // Deposit only при расходе -> Ineligible transaction
	}
	return ""
}
