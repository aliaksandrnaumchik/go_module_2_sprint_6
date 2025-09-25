package service

import (
	"bytes"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

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
		utf8String := convertWindows1251ToUTF8(input)
		return morse.ToMorse(utf8String), nil
	}

	return "", fmt.Errorf("пока не умею работать с этой кодировкой itf8")
}

func isMorse(input string) bool {
	for _, char := range input {
		if char != '.' && char != '-' && char != ' ' {
			return false
		}
	}
	return true
}

func isPlainText(input string) bool {
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

func isWindows1251(input string) bool {
	windows1251Bytes := []byte(input)
	decoder := charmap.Windows1251.NewDecoder()
	reader := transform.NewReader(bytes.NewReader(windows1251Bytes), decoder)
	_, err := io.Copy(io.Discard, reader)
	return err == nil
}

func convertWindows1251ToUTF8(input string) string {
	windows1251Bytes := []byte(input)
	decoder := charmap.Windows1251.NewDecoder()
	utf8String, err := decoder.Bytes(windows1251Bytes)
	if err != nil {
		fmt.Printf("Ошибка декодирования: %v", err)
	}
	return string(utf8String)
}
