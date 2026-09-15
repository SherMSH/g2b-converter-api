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
	want := map[string]string{
		"0": "50", "2": "40", "3": "41", "4": "58",
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
