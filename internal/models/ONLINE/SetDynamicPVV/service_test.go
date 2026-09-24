package setdynamicpvv

import "testing"

// В PINBlock ждём открытый PIN: блок под рабочим ключом расшифровать нечем
func TestIsClearPIN(t *testing.T) {
	valid := []string{"1234", "0000", "123456", "123456789012"}
	for _, pin := range valid {
		if !isClearPIN(pin) {
			t.Errorf("%q - допустимый PIN", pin)
		}
	}

	invalid := map[string]string{
		"":                 "пустое значение",
		"123":              "короче четырёх цифр",
		"1234567890123":    "длиннее двенадцати",
		"12a4":             "не только цифры",
		"D58221C0227900B5": "шифрованный PIN-блок",
		" 1234":            "с пробелом",
	}
	for pin, why := range invalid {
		if isClearPIN(pin) {
			t.Errorf("%q (%s) не должен приниматься", pin, why)
		}
	}
}
