package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

const INDEX_HTML = "index.html"

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

	parentDir := filepath.Dir(cwd)
	filePath := filepath.Join(parentDir, INDEX_HTML)

	r.Header.Add("Content-Type", "text/html")
	http.ServeFile(w, r, filePath)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // 32 MB limit
	if err != nil {
		log.Printf("Ошибка при парсинге формы: %v", err)
		http.Error(w, "Ошибка загрузки файла", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Printf("Ошибка при получении файла: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusBadRequest)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Ошибка при чтении файла: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(convertedContent))
}
