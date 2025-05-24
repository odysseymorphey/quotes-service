package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	handlers2 "github.com/odysseymorphey/quotes-service/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/odysseymorphey/quotes-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBaseHandler_AddQuote(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		mockError    error
		expectedCode int
		expectedBody string
	}{
		{
			name: "successful add",
			requestBody: models.Quote{
				Author: "Test Author",
				Quote:  "Test Quote",
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid json",
			requestBody:  "invalid json",
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid request\n",
		},
		{
			name: "repository error",
			requestBody: models.Quote{
				Author: "Test Author",
				Quote:  "Test Quote",
			},
			mockError:    errors.New("database error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Internal server error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			handler := &handlers2.BaseHandler{Repo: mockRepo}

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/quotes", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			if tt.mockError != nil {
				mockRepo.On("AddQuote", mock.Anything, mock.AnythingOfType("models.Quote")).
					Return(tt.mockError)
			} else if tt.requestBody != "invalid json" {
				mockRepo.On("AddQuote", mock.Anything, mock.AnythingOfType("models.Quote")).
					Return(nil)
			}

			handler.AddQuote(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)

			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}

			if tt.mockError == nil && tt.requestBody != "invalid json" {
				mockRepo.AssertCalled(t, "AddQuote", mock.Anything, mock.MatchedBy(func(q models.Quote) bool {
					return q.Author == tt.requestBody.(models.Quote).Author &&
						q.Quote == tt.requestBody.(models.Quote).Quote
				}))
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBaseHandler_DeleteQuote(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockError      error
		expectedCode   int
		expectedBody   string
		expectedLogMsg string
	}{
		{
			name:         "successful delete",
			id:           "valid-id",
			expectedCode: http.StatusOK,
		},
		{
			name:         "repository error",
			id:           "error-id",
			mockError:    errors.New("database failure"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Internal server error\n",
		},
		{
			name:         "not found",
			id:           "non-existent-id",
			mockError:    fmt.Errorf("postgres.DeleteQuote: quote not found"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Internal server error\n",
		},
		{
			name:         "empty id",
			id:           "",
			mockError:    fmt.Errorf("postgres.DeleteQuote: quote not found"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			handler := &handlers2.BaseHandler{Repo: mockRepo}

			req := httptest.NewRequest("DELETE", "/quotes/"+tt.id, nil)
			if tt.id != "" {
				req.SetPathValue("id", tt.id)
			}

			rr := httptest.NewRecorder()

			if tt.id != "" || tt.mockError != nil {
				mockRepo.On("DeleteQuote", mock.Anything, tt.id).
					Return(tt.mockError).
					Once()
			}

			handler.DeleteQuote(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}

			if tt.id != "" || tt.mockError != nil {
				mockRepo.AssertCalled(t, "DeleteQuote", mock.Anything, tt.id)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBaseHandler_GetQuotes(t *testing.T) {
	mockQuotes := []models.Quote{
		{Id: "1", Author: "Author1", Quote: "Quote1"},
		{Id: "2", Author: "Author2", Quote: "Quote2"},
	}

	tests := []struct {
		name           string
		queryParams    map[string]string
		mockQuotes     []models.Quote
		mockError      error
		expectedCode   int
		expectedBody   string
		expectedHeader string
	}{
		{
			name:           "get all quotes successfully",
			mockQuotes:     mockQuotes,
			expectedCode:   http.StatusOK,
			expectedBody:   `[{"id":"1","author":"Author1","quote":"Quote1"},{"id":"2","author":"Author2","quote":"Quote2"}]` + "\n",
			expectedHeader: "application/json",
		},
		{
			name:           "get quotes by author",
			queryParams:    map[string]string{"author": "Author1"},
			mockQuotes:     []models.Quote{mockQuotes[0]},
			expectedCode:   http.StatusOK,
			expectedBody:   `[{"id":"1","author":"Author1","quote":"Quote1"}]` + "\n",
			expectedHeader: "application/json",
		},
		{
			name:         "error getting all quotes",
			mockError:    errors.New("database error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Internal server error\n",
		},
		{
			name:         "error getting quotes by author",
			queryParams:  map[string]string{"author": "Unknown"},
			mockError:    errors.New("not found"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Internal server error\n",
		},
		{
			name:           "empty quotes list",
			mockQuotes:     []models.Quote{},
			expectedCode:   http.StatusOK,
			expectedBody:   "[]\n",
			expectedHeader: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			handler := &handlers2.BaseHandler{Repo: mockRepo}

			req := httptest.NewRequest("GET", "/quotes", nil)
			q := req.URL.Query()
			for k, v := range tt.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			rr := httptest.NewRecorder()

			// Настройка моков с учетом ошибок
			if author := tt.queryParams["author"]; author != "" {
				mockRepo.On("GetQuotesByAuthor", mock.Anything, author).
					Return(tt.mockQuotes, tt.mockError)
			} else {
				mockRepo.On("GetQuotes", mock.Anything).
					Return(tt.mockQuotes, tt.mockError)
			}

			handler.GetQuotes(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBaseHandler_GetRandomQuote(t *testing.T) {
	mockQuote := &models.Quote{
		Id:     "1",
		Author: "Test Author",
		Quote:  "Test Quote",
	}

	tests := []struct {
		name           string
		mockQuote      *models.Quote
		mockError      error
		expectedCode   int
		expectedBody   string
		expectedHeader string
	}{
		{
			name:           "successful get random quote",
			mockQuote:      mockQuote,
			expectedCode:   http.StatusOK,
			expectedBody:   `{"id":"1","author":"Test Author","quote":"Test Quote"}` + "\n",
			expectedHeader: "application/json",
		},
		{
			name:           "repository error",
			mockError:      errors.New("database error"),
			expectedCode:   http.StatusInternalServerError,
			expectedBody:   "Internal server error\n",
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "no quotes found",
			mockError:      sql.ErrNoRows,
			expectedCode:   http.StatusNotFound,
			expectedBody:   "No quotes found\n",
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "nil quote without error",
			mockQuote:      nil,
			expectedCode:   http.StatusNotFound,
			expectedBody:   "No quotes found\n",
			expectedHeader: "text/plain; charset=utf-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			handler := &handlers2.BaseHandler{Repo: mockRepo}

			mockRepo.On("GetRandomQuote", mock.Anything).
				Return(tt.mockQuote, tt.mockError)

			req := httptest.NewRequest("GET", "/quotes/random", nil)
			rr := httptest.NewRecorder()

			handler.GetRandomQuote(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			assert.Equal(t, tt.expectedHeader, rr.Header().Get("Content-Type"))

			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}
		})
	}
}
