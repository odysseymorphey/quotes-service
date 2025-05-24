package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func (h *BaseHandler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote, err := h.Repo.GetRandomQuote(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No quotes found", http.StatusNotFound)
			return
		}
		log.Printf("Can't get quote: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if quote == nil {
		http.Error(w, "No quotes found", http.StatusNotFound)
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
