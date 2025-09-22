package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "app: ", log.LstdFlags|log.Lshortfile)

	srv := server.NewServer(logger)

	if err := srv.Start(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
