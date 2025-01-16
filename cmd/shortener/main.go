package main

import (
	"log"
	"net/http"

	"github.com/Arth19/url-shortener-go/internal/app"
	"github.com/Arth19/url-shortener-go/internal/storage"
	"github.com/gorilla/mux"
)

func main() {
	err := storage.InitPostgres()
	if err != nil {
		log.Fatalf("Falha ao inicializar Postgres: %v", err)
	}

	router := mux.NewRouter()
	app.SetupRoutes(router)

	log.Println("Servidor iniciado na porta :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}
}
