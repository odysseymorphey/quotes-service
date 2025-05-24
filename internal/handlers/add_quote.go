package handlers

import (
	"encoding/json"
	"github.com/odysseymorphey/quotes-service/internal/models"
	"log"
	"net/http"
)

func (h *BaseHandler) AddQuote(w http.ResponseWriter, r *http.Request) {
	var quote models.Quote
	if err := json.NewDecoder(r.Body).Decode(&quote); err != nil {
		log.Printf("Failed request body decoding: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.repo.AddQuote(r.Context(), quote); err != nil {
		log.Printf("Failed to add quote: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
