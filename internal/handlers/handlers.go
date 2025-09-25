package handlers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Ошибка при получении текущей директории: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(cwd, "../index.html")

	templateFile, err := template.ParseFiles(filePath)
	if err != nil {
		log.Printf("Ошибка при парсинге шаблона: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	if err := templateFile.Execute(w, nil); err != nil {
		log.Printf("Ошибка при выполнении шаблона: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // 32 MB limit
	if err != nil {
		log.Printf("Ошибка при парсинге формы: %v", err)
		http.Error(w, "Ошибка загрузки файла", http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Printf("Ошибка при получении файла: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Ошибка при чтении файла: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	convertedContent, err := service.AutoConvert(string(content))
	if err != nil {
		log.Printf("Ошибка при конвертации: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	timestamp := time.Now().UTC().String()
	ext := filepath.Ext(handler.Filename)
	newFilename := fmt.Sprintf("converted_%s%s", timestamp, ext)

	f, err := os.Create(newFilename)
	if err != nil {
		log.Printf("Ошибка при создании файла: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = f.WriteString(convertedContent)
	if err != nil {
		log.Printf("Ошибка при записи в файл: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	err = f.Sync()
	if err != nil {
		log.Printf("Ошибка при синхронизации файла: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(convertedContent))
}
