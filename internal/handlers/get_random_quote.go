package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *BaseHandler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote, err := h.repo.GetRandomQuote(r.Context())
	if err != nil {
		log.Printf("Can't get quote: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(quote); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
