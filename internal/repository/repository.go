package repository

import (
	"context"
	"github.com/odysseymorphey/quotes-service/internal/models"
)

type Repository interface {
	AddQuote(ctx context.Context, q models.Quote) error
	GetQuotes(ctx context.Context) ([]models.Quote, error)
	GetQuotesByAuthor(ctx context.Context, author string) ([]models.Quote, error)
	GetRandomQuote(ctx context.Context) (*models.Quote, error)
	DeleteQuote(ctx context.Context, id string) error
	Close() error
}
