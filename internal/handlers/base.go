package handlers

import "github.com/odysseymorphey/quotes-service/internal/repository"

type BaseHandler struct {
	Repo repository.Repository
}

func New(r repository.Repository) *BaseHandler {
	return &BaseHandler{
		Repo: r,
	}
}
