package service

import (
	"bytes"
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

// Полная таблица соответствий Windows-1251 -> UTF-8 для кириллицы
var windows1251Map = map[byte]string{
	// Заглавные буквы
	0xC0: "А", 0xC1: "Б", 0xC2: "В", 0xC3: "Г", 0xC4: "Д",
	0xC5: "Е", 0xC6: "Ж", 0xC7: "З", 0xC8: "И", 0xC9: "Й",
	0xCA: "К", 0xCB: "Л", 0xCC: "М", 0xCD: "Н", 0xCE: "О",
	0xCF: "П", 0xD0: "Р", 0xD1: "С", 0xD2: "Т", 0xD3: "У",
	0xD4: "Ф", 0xD5: "Х", 0xD6: "Ц", 0xD7: "Ч", 0xD8: "Ш",
	0xD9: "Щ", 0xDA: "Ъ", 0xDB: "Ы", 0xDC: "Ь", 0xDD: "Э",
	0xDE: "Ю", 0xDF: "Я",

	// Строчные буквы
	0xE0: "а", 0xE1: "б", 0xE2: "в", 0xE3: "г", 0xE4: "д",
	0xE5: "е", 0xE6: "ж", 0xE7: "з", 0xE8: "и", 0xE9: "й",
	0xEA: "к", 0xEB: "л", 0xEC: "м", 0xED: "н", 0xEE: "о",
	0xEF: "п", 0xF0: "р", 0xF1: "с", 0xF2: "т", 0xF3: "у",
	0xF4: "ф", 0xF5: "х", 0xF6: "ц", 0xF7: "ч", 0xF8: "ш",
	0xF9: "щ", 0xFA: "ъ", 0xFB: "ы", 0xFC: "ь", 0xFD: "э",
	0xFE: "ю", 0xFF: "я",
}

func AutoConvert(input string) (string, error) {
	if isMorse(input) {
		text := morse.ToText(input)
		return text, nil
	} else if isPlainText(input) {
		morse := morse.ToMorse(input)
		return morse, nil
	} else if !isUtf8(input) {
		return utf8Convert(input)
	}
	return "", fmt.Errorf("ошибка конвертации в Морзе: %s", input)
}

func utf8Convert(input string) (string, error) {
	if isWindows1251(input) {
		windows1251Bytes := []byte(input)
		utf8String := convertWindows1251ToUTF8(windows1251Bytes)
		return morse.ToMorse(utf8String), nil
	}

	return "", fmt.Errorf("пока не умею работать с этой кодировкой")
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

func isUtf8(text string) bool {
	return utf8.ValidString(text)
}

func isWindows1251(s string) bool {
	for _, b := range []byte(s) {
		// Проверяем, что все байты находятся в допустимом диапазоне Windows-1251
		if b < 0 || b > 255 {
			return false
		}
	}
	return true
}

// Функция преобразования
func convertWindows1251ToUTF8(input []byte) string {
	var buf bytes.Buffer
	for _, b := range input {
		if utf8Char, ok := windows1251Map[b]; ok {
			buf.WriteString(utf8Char)
		} else {
			// Если символ не найден, добавляем как есть
			buf.WriteByte(b)
		}
	}
	return buf.String()
}
