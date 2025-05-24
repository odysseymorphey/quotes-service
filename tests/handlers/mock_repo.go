package handlers

import (
	"context"
	"github.com/odysseymorphey/quotes-service/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) AddQuote(ctx context.Context, q models.Quote) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

func (m *MockRepository) GetQuotes(ctx context.Context) ([]models.Quote, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Quote), args.Error(1)
}

func (m *MockRepository) GetQuotesByAuthor(ctx context.Context, author string) ([]models.Quote, error) {
	args := m.Called(ctx, author)
	return args.Get(0).([]models.Quote), args.Error(1)
}

func (m *MockRepository) GetRandomQuote(ctx context.Context) (*models.Quote, error) {
	args := m.Called(ctx)
	return args.Get(0).(*models.Quote), args.Error(1)
}

func (m *MockRepository) DeleteQuote(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) Close() error {
	args := m.Called()
	return args.Error(0)
}
