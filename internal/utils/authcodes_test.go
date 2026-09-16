package utils

import "testing"

func TestAuthRespCode(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		rspcode string
		want    string
	}{
		{"одобрено", "0", "00", TwoApproved},
		{"одобрено частично", "0", "02", TwoApprovedPartial},
		{"нулевая сумма - проверка карты", "0", "85", TwoApproved},
		{"недостаточно средств", "1", "16", "59"},
		{"истёк срок действия", "1", "01", "51"},
		{"карта ограничена", "1", "04", "58"},
		{"неверный PIN", "1", "17", "53"},
		{"лимит вводов PIN достигнут", "1", "06", "62"},
		{"лимит вводов PIN был достигнут ранее", "1", "93", "83"},
		{"ошибка проверки CVV", "1", "85", "81"},
		{"подозрение в мошенничестве", "1", "02", "75"},
		{"изъятие карты", "2", "00", "50"},
		{"системная ошибка процессинга", "9", "09", TwoSystemError},
		{"неизвестный отказ", "1", "99999", TwoExternalDecline},
		{"неизвестное одобрение", "0", "99999", TwoApproved},
		{"пустой ответ", "", "", TwoNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AuthRespCode(tt.code, tt.rspcode); got != tt.want {
				t.Errorf("AuthRespCode(%q, %q) = %q, want %q", tt.code, tt.rspcode, got, tt.want)
			}
		})
	}
}

func TestIsApprovedAndRetainCard(t *testing.T) {
	if !IsApproved("0") {
		t.Error("код 0 - одобрение")
	}
	if IsApproved("1") {
		t.Error("код 1 - отказ, не одобрение")
	}
	if RetainCard("2") != "1" {
		t.Error("код 2 требует изъятия карты")
	}
	if RetainCard("1") != "0" {
		t.Error("обычный отказ не требует изъятия карты")
	}
}

// Коды по статусу карты из таблицы «Статусы карты в TWO и коды ответов
// авторизатора». Статусы 1, 5 и 6 операции не запрещают.
func TestCardStatusRespCodes(t *testing.T) {
	// Статус 4 зависит от направления операции и проверяется отдельно
	want := map[string]string{
		"0": "50", "2": "40", "3": "41",
		"8": "75", "9": "50", "10": "71", "12": "50", "15": "51",
	}
	for status, code := range want {
		if got := CardStatusRespCodes[status]; got != code {
			t.Errorf("статус карты %s -> %q, want %q", status, got, code)
		}
	}
	for _, status := range []string{"1", "5", "6"} {
		if _, blocked := CardStatusRespCodes[status]; blocked {
			t.Errorf("статус карты %s не должен запрещать операции", status)
		}
	}
}

// «Только приход» отклоняет расход и пропускает зачисление.
func TestAccountStatusRespCode(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		isDebit bool
		want    string
	}{
		{"неактивный счёт", "0", true, "56"},
		{"закрытый счёт", "9", false, "56"},
		{"только просмотр", "5", false, "55"},
		{"только приход, расход", "2", true, "55"},
		{"только приход, зачисление", "2", false, ""},
		{"только приход первичный, расход", "4", true, "55"},
		{"только приход первичный, зачисление", "4", false, ""},
		{"открытый счёт", "1", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AccountStatusRespCode(tt.status, tt.isDebit); got != tt.want {
				t.Errorf("AccountStatusRespCode(%q, %v) = %q, want %q", tt.status, tt.isDebit, got, tt.want)
			}
		})
	}
}

// Семейство 9 - это «ошибки обработки», но не все они системные: часть несёт
// обычную причину отказа и не должна схлопываться в 54.
func TestAuthRespCodeErrorFamily(t *testing.T) {
	tests := []struct {
		rspcode string
		want    string
		name    string
	}{
		{"09", TwoSystemError, "системный сбой"},
		{"58", "59", "недостаточно средств"},
		{"59", "53", "неверный PIN"},
		{"56", "52", "неверный номер карты"},
		{"51", "51", "истёк срок действия"},
		{"14", "15", "оригинал не найден"},
		{"12", "72", "эмитент недоступен"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AuthRespCode("9", tt.rspcode); got != tt.want {
				t.Errorf("AuthRespCode(9, %s) = %q, want %q", tt.rspcode, got, tt.want)
			}
		})
	}
}

// Служебные коды уровня сервиса тоже должны доходить до партнёра осмысленными
func TestAuthRespCodeServiceLevel(t *testing.T) {
	cases := map[string]string{"D/00": "52", "D/04": "15", "D/07": "56", "C/43": "74"}
	for key, want := range cases {
		code, rspcode := key[:1], key[2:]
		if got := AuthRespCode(code, rspcode); got != want {
			t.Errorf("AuthRespCode(%s, %s) = %q, want %q", code, rspcode, got, want)
		}
	}
}

// Изъятие карты по утере и краже
func TestAuthRespCodePickupLostStolen(t *testing.T) {
	if got := AuthRespCode("2", "08"); got != "40" {
		t.Errorf("утерянная карта -> %q, want 40", got)
	}
	if got := AuthRespCode("2", "09"); got != "41" {
		t.Errorf("украденная карта -> %q, want 41", got)
	}
}

// Статус 4 (Restricted) разрешает зачисление и запрещает расход
func TestCardStatusRespCodeRestricted(t *testing.T) {
	if got := CardStatusRespCode("4", true); got != "58" {
		t.Errorf("расход по ограниченной карте -> %q, want 58", got)
	}
	if got := CardStatusRespCode("4", false); got != "" {
		t.Errorf("зачисление на ограниченную карту не запрещено, получено %q", got)
	}
	// Остальные статусы направление не учитывают
	if got := CardStatusRespCode("8", false); got != "75" {
		t.Errorf("скомпрометированная карта -> %q, want 75", got)
	}
	if got := CardStatusRespCode("1", true); got != "" {
		t.Errorf("открытая карта не запрещает операции, получено %q", got)
	}
}

// D8 statcode 16 - «использована мошеннически», по таблице партнёра это
// статус 8 (скомпрометирована), а не 4 (ограничена)
func TestCardStatusesFraudulent(t *testing.T) {
	if got := CardStatuses["16"]; got != "8" {
		t.Errorf("statCode 16 -> %q, want 8", got)
	}
	if got := CardStatuses["08"]; got != "8" {
		t.Errorf("statCode 08 -> %q, want 8", got)
	}
	if got := CardStatuses["10"]; got != "4" {
		t.Errorf("statCode 10 -> %q, want 4", got)
	}
}

// Справочник статусов процессинга: каждый статус карты D8 должен переводиться
// в статус TWO, иначе в ответ партнёру уйдёт заглушка "unknown".
func TestCardStatusesCoversD8Dictionary(t *testing.T) {
	d8 := []string{"00", "01", "02", "03", "04", "05", "06", "08", "10", "11",
		"12", "13", "14", "15", "16", "17", "20", "21", "22", "23", "24", "25"}
	for _, code := range d8 {
		if two, ok := CardStatuses[code]; !ok || two == "" || two == "unknown" {
			t.Errorf("statCode %s не переводится в статус TWO: %q", code, two)
		}
	}
}
