package server

import (
	"context"
	"log"
	"net/http"

	"github.com/odysseymorphey/quotes-service/internal/handlers"
	"github.com/odysseymorphey/quotes-service/internal/repository"
)

type Server struct {
	srv  *http.Server
	repo repository.Repository
}

func New(r repository.Repository) *Server {
	m := http.NewServeMux()
	h := handlers.New(r)

	registerRoutes(m, h)

	return &Server{
		srv: &http.Server{
			Addr:    ":8080",
			Handler: m,
		},
		repo: r,
	}
}

func registerRoutes(mux *http.ServeMux, h *handlers.BaseHandler) {
	mux.HandleFunc("POST /quotes", h.AddQuote)

	mux.HandleFunc("GET /quotes", h.GetQuotes)
	mux.HandleFunc("GET /quotes/random", h.GetRandomQuote)

	mux.HandleFunc("DELETE /quotes/{id}", h.DeleteQuote)
}

func (s *Server) Run() {
	if err := s.srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func (s *Server) Stop() {
	log.Println(s.repo.Close())
	log.Println(s.srv.Shutdown(context.Background()))
}
