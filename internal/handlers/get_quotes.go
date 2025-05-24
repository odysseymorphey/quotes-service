package handlers

import (
	"encoding/json"
	"github.com/odysseymorphey/quotes-service/internal/models"
	"log"
	"net/http"
)

func (h *BaseHandler) GetQuotes(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")

	var quotes []models.Quote
	var err error

	if author != "" {
		quotes, err = h.Repo.GetQuotesByAuthor(r.Context(), author)
	} else {
		quotes, err = h.Repo.GetQuotes(r.Context())
	}

	if err != nil {
		log.Printf("Can't get quotes: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(quotes); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
