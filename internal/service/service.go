package service

import (
	"fmt"
	"unicode"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func AutoConvert(input string) (string, error) {
	if isMorse(input) {
		text := morse.ToText(input)
		return text, nil
	} else if isPlainText(input) {
		morse := morse.ToMorse(input)
		return morse, nil
	} else {
		return "", fmt.Errorf("ошибка конвертации в Морзе: %s", input)
	}
}

func isMorse(input string) bool {
	// Разрешаем только точки, тире и пробелы
	for _, char := range input {
		if char != '.' && char != '-' && char != ' ' {
			return false
		}
	}
	return true
}

func isPlainText(input string) bool {
	// Проверяем, содержит ли строка буквенные символы
	for _, char := range input {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			return true
		}
	}
	return false
}
