package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Arth19/url-shortener-go/internal/shortener"
	"github.com/Arth19/url-shortener-go/internal/storage"
	"github.com/gorilla/mux"
)

type RequestBody struct {
	URL string `json:"url"`
}

type ResponseBody struct {
	ShortCode string `json:"short_code"`
}

type StatsResponse struct {
	ShortCode  string `json:"short_code"`
	Original   string `json:"original"`
	ClickCount uint   `json:"click_count"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	urlData, err := storage.GetURL(shortCode)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	resp := StatsResponse{
		ShortCode:  urlData.ShortCode,
		Original:   urlData.Original,
		ClickCount: urlData.ClickCount,
		CreatedAt:  urlData.CreatedAt.Format("02/01/2006 15:04:05"),
		UpdatedAt:  urlData.UpdatedAt.Format("02/01/2006 15:04:05"),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("Erro ao codificar JSON:", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
}

func EncurtarURLHandler(w http.ResponseWriter, r *http.Request) {
	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	if body.URL == "" {
		http.Error(w, "URL não pode ser vazia", http.StatusBadRequest)
		return
	}

	shortCode, err := shortener.GenerateShortCode()
	if err != nil {
		log.Println("Erro ao gerar short code:", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	_, err = storage.SaveURL(shortCode, body.URL)
	if err != nil {
		log.Println("Erro ao salvar no banco:", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	resp := ResponseBody{ShortCode: shortCode}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("Erro ao codificar JSON:", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
}

func RedirecionarHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := mux.Vars(r)["shortCode"]
	originalURL, err := storage.GetURL(shortCode)
	if err != nil {
		log.Println("Erro ao obter URL original:", err)
		http.NotFound(w, r)
		return
	}

	err = storage.IncrementClickCount(shortCode)
	if err != nil {
		log.Println("Erro ao incrementar contagem de cliques:", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	log.Printf("Redirecionando para %s", originalURL.Original)
	http.Redirect(w, r, originalURL.Original, http.StatusFound)
}
