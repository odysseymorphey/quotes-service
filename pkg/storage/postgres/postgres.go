package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/odysseymorphey/quotes-service/internal/models"
)

type Database struct {
	db *sql.DB
}

func New(dsn string) (*Database, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("can't ping database")
	}

	return &Database{
		db: db,
	}, nil
}

func (d *Database) AddQuote(ctx context.Context, q models.Quote) error {
	const op = "postgres.AddQuote"

	query := `INSERT INTO quotes(author, quote) VALUES ($1, $2)`

	_, err := d.db.ExecContext(ctx, query, q.Author, q.Quote)
	if err != nil {
		return fmt.Errorf("%s: failed to execute query: %v", op, err)
	}

	return nil
}

func (d *Database) GetQuotes(ctx context.Context) ([]models.Quote, error) {
	const op = "postgres.GetQuotes"

	query := `SELECT * FROM quotes`

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query: %v", op, err)
	}

	var quotes []models.Quote
	for rows.Next() {
		var quote models.Quote
		if err := rows.Scan(&quote.Id, &quote.Author, &quote.Quote); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %v", op, err)
		}

		quotes = append(quotes, quote)
	}

	return quotes, nil
}

func (d *Database) GetQuotesByAuthor(ctx context.Context, author string) ([]models.Quote, error) {
	const op = "postgres.GetQuotesByAuthor"

	query := `SELECT * FROM quotes WHERE author = $1`

	rows, err := d.db.QueryContext(ctx, query, author)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query: %v", op, err)
	}

	var quotes []models.Quote
	for rows.Next() {
		var quote models.Quote
		if err := rows.Scan(&quote.Id, &quote.Author, &quote.Quote); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %v", op, err)
		}

		quotes = append(quotes, quote)
	}

	return quotes, nil
}

func (d *Database) GetRandomQuote(ctx context.Context) (*models.Quote, error) {
	const op = "postgres.GetRandomQuote"

	query := `SELECT * FROM quotes ORDER BY random() LIMIT 1`

	row := d.db.QueryRowContext(ctx, query)

	var quote models.Quote
	err := row.Scan(&quote.Id, &quote.Author, &quote.Quote)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to scan row: %v", op, err)
	}

	return &quote, nil
}
func (d *Database) DeleteQuote(ctx context.Context, id string) error {
	const op = "postgres.DeleteQuote"

	query := `DELETE FROM quotes WHERE id = $1`

	res, err := d.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: quote not found", op)
	}

	return nil
}

func (d *Database) Close() error {
	if err := d.db.Close(); err != nil {
		return fmt.Errorf("failed to close database: %v", err)
	}

	return nil
}
