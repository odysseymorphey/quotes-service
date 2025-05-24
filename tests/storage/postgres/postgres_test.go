package postgres

import (
	"context"
	"database/sql"
	"errors"
	postgres2 "github.com/odysseymorphey/quotes-service/pkg/storage/postgres"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odysseymorphey/quotes-service/internal/models"
	"github.com/stretchr/testify/assert"
)

func NewMock() (*postgres2.Database, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	return &postgres2.Database{db: db}, mock
}

func TestAddQuote(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	q := models.Quote{
		Author: "Test Author",
		Quote:  "Test Quote",
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO quotes").
			WithArgs(q.Author, q.Quote).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := db.AddQuote(context.Background(), q)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO quotes").
			WithArgs(q.Author, q.Quote).
			WillReturnError(errors.New("connection failed"))

		err := db.AddQuote(context.Background(), q)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to execute query")
	})
}

func TestGetQuotes(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "author", "quote"}).
			AddRow(1, "Author1", "Quote1").
			AddRow(2, "Author2", "Quote2")

		mock.ExpectQuery("SELECT \\* FROM quotes").WillReturnRows(rows)

		quotes, err := db.GetQuotes(context.Background())
		assert.NoError(t, err)
		assert.Len(t, quotes, 2)
	})

	t.Run("Empty", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "author", "quote"})
		mock.ExpectQuery("SELECT \\* FROM quotes").WillReturnRows(rows)

		quotes, err := db.GetQuotes(context.Background())
		assert.NoError(t, err)
		assert.Empty(t, quotes)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM quotes").
			WillReturnError(errors.New("query error"))

		_, err := db.GetQuotes(context.Background())
		assert.Error(t, err)
	})
}

func TestGetQuotesByAuthor(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	author := "TestAuthor"

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "author", "quote"}).
			AddRow(1, author, "Quote1").
			AddRow(2, author, "Quote2")

		mock.ExpectQuery("SELECT \\* FROM quotes WHERE author = ?").
			WithArgs(author).
			WillReturnRows(rows)

		quotes, err := db.GetQuotesByAuthor(context.Background(), author)
		assert.NoError(t, err)
		assert.Len(t, quotes, 2)
		assert.Equal(t, author, quotes[0].Author)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM quotes WHERE author = ?").
			WithArgs(author).
			WillReturnError(sql.ErrNoRows)

		_, err := db.GetQuotesByAuthor(context.Background(), author)
		assert.Error(t, err)
	})
}

func TestGetRandomQuote(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	t.Run("Success", func(t *testing.T) {
		row := sqlmock.NewRows([]string{"id", "author", "quote"}).
			AddRow(1, "Author", "Quote")

		mock.ExpectQuery("^SELECT \\* FROM quotes ORDER BY random\\(\\) LIMIT 1$").
			WillReturnRows(row)

		quote, err := db.GetRandomQuote(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, quote)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT * FROM quotes ORDER BY random() LIMIT 1").
			WillReturnError(sql.ErrNoRows)

		_, err := db.GetRandomQuote(context.Background())
		assert.Error(t, err)
	})
}

func TestDeleteQuote(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	id := "123"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM quotes WHERE id = ?").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := db.DeleteQuote(context.Background(), id)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM quotes WHERE id = ?").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := db.DeleteQuote(context.Background(), id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "quote not found")
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM quotes WHERE id = ?").
			WithArgs(id).
			WillReturnError(errors.New("db error"))

		err := db.DeleteQuote(context.Background(), id)
		assert.Error(t, err)
	})
}

func TestClose(t *testing.T) {
	db, mock := NewMock()

	mock.ExpectClose()
	err := db.Close()
	assert.NoError(t, err)
}

func TestContextHandling(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	mock.ExpectExec("INSERT INTO quotes").WillReturnResult(sqlmock.NewResult(1, 1))
	<-ctx.Done()
	err := db.AddQuote(ctx, models.Quote{})
	assert.Error(t, err)
}
