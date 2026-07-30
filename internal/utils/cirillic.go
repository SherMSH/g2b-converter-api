package utils

import "strings"

// TajikToCyrillic преобразует таджикские национальные символы в строке на кириллические
func TajikToCyrillic(text string) string {
	// Таблица соответствия таджикских символов -> кириллица
	translitMap := map[string]string{
		// Буквы
		"Ғ": "Г", "ғ": "г",
		"Ӣ": "И", "ӣ": "и",
		"Қ": "К", "қ": "к",
		"Ӯ": "У", "ӯ": "у",
		"Ҳ": "Х", "ҳ": "х",
		"Ҷ": "Ч", "ҷ": "ч",
	}

	result := text

	keys := []string{"Ғ", "Ӣ", "Қ", "Ӯ", "Ҳ", "Ҷ", "ғ", "ӣ", "қ", "ӯ", "ҳ", "ҷ"}
	for _, key := range keys {
		result = strings.ReplaceAll(result, key, translitMap[key])
	}
	return result
}
