package runes

import "unicode"

func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}

	// Создаём срез рун для корректной обработки UTF-8 символов
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
