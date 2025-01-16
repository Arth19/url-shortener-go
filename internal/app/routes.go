package app

import (
	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/encurtar", EncurtarURLHandler).Methods("POST")
	r.HandleFunc("/{shortCode}", RedirecionarHandler).Methods("GET")
	r.HandleFunc("/stats/{shortCode}", StatsHandler).Methods("GET")
}
