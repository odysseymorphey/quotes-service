package handlers

import (
	"log"
	"net/http"
)

func (h *BaseHandler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.repo.DeleteQuote(r.Context(), id); err != nil {
		log.Printf("Failed to delete quote: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
