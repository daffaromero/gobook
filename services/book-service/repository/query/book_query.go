package query

import (
	"context"
	"errors"
	"fmt"
	"time"

	api "github.com/daffaromero/gobook/protobuf/api"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookQuery interface {
	GetBook(ctx context.Context, req *api.GetBookRequest) (*api.GetBookResponse, error)
	ListBooks(ctx context.Context, req *api.ListBooksRequest) (*api.ListBooksResponse, error)
	CreateBook(ctx context.Context, tx pgx.Tx, req *api.CreateBookRequest) (*api.CreateBookResponse, error)
	UpdateBook(ctx context.Context, tx pgx.Tx, req *api.UpdateBookRequest) (*api.UpdateBookResponse, error)
	DeleteBook(ctx context.Context, tx pgx.Tx, id *api.DeleteBookRequest) (*api.DeleteBookResponse, error)
}

type bookQuery struct {
	db *pgxpool.Pool
}

func NewBookQuery(db *pgxpool.Pool) *bookQuery {
	return &bookQuery{
		db: db,
	}
}

func (q *bookQuery) GetBook(ctx context.Context, req *api.GetBookRequest) (*api.GetBookResponse, error) {
	if req == nil || req.BookId == "" {
		return nil, errors.New("book ID cannot be empty")
	}
	query := `SELECT id, title, author, category_id, description, published_date FROM books WHERE id = $1 AND deleted_at IS NULL`

	row := q.db.QueryRow(ctx, query, req.BookId)

	var book api.Book
	err := row.Scan(&book.Id, &book.Title, &book.Author, &book.CategoryId, &book.Description, &book.PublishedDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("book with ID %s not found", req.BookId)
		}
		return nil, fmt.Errorf("failed to scan book: %w", err)
	}

	return &api.GetBookResponse{
		Book: &book,
	}, nil
}

func (q *bookQuery) ListBooks(ctx context.Context, req *api.ListBooksRequest) (*api.ListBooksResponse, error) {
	query := `SELECT id, title, author, category_id, description, published_date FROM books WHERE deleted_at IS NULL`

	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query books: %w", err)
	}
	defer rows.Close()

	var books []*api.Book
	for rows.Next() {
		var book api.Book
		err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.CategoryId, &book.Description, &book.PublishedDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan book: %w", err)
		}
		books = append(books, &book)
	}

	return &api.ListBooksResponse{
		Books: books,
	}, nil
}

func (q *bookQuery) CreateBook(ctx context.Context, tx pgx.Tx, req *api.CreateBookRequest) (*api.CreateBookResponse, error) {
	if req == nil || req.Book == nil {
		return nil, errors.New("book cannot be empty")
	}
	query := `INSERT INTO books (id, title, author, category_id, description, published_date) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, title, author, category_id, description, published_date`

	createdAt := req.Book.CreatedAt.AsTime()
	updatedAt := req.Book.UpdatedAt.AsTime()

	var createdBook api.Book

	err := tx.QueryRow(ctx, query, req.Book.Id, req.Book.Title, req.Book.Author, req.Book.CategoryId, req.Book.Description, req.Book.PublishedDate, createdAt, updatedAt).Scan(&createdBook.Id, &createdBook.Title, &createdBook.Author, &createdBook.CategoryId, &createdBook.Description, &createdBook.PublishedDate)
	if err != nil {
		return nil, fmt.Errorf("failed to insert book: %w", err)
	}

	return &api.CreateBookResponse{
		Book: &createdBook,
	}, nil
}

func (q *bookQuery) UpdateBook(ctx context.Context, tx pgx.Tx, req *api.UpdateBookRequest) (*api.UpdateBookResponse, error) {
	if req == nil || req.Book == nil {
		return nil, errors.New("book cannot be empty")
	}
	query := `UPDATE books 
		SET 
			title = COALESCE($2, name), 
			author = COALESCE($3, author), 
			category_id = COALESCE($4, category_id), 
			description = COALESCE($5, description), 
			published_date = COALESCE($6, published_date),
			updated_at = $7
		WHERE id = $1 AND deleted_at IS NULL RETURNING id, title, author, category_id, description, published_date`

	updatedAt := time.Now()
	if req.Book.UpdatedAt != nil {
		updatedAt = req.Book.UpdatedAt.AsTime()
	}

	var updatedBook api.Book

	err := tx.QueryRow(ctx, query, req.Book.Id, req.Book.Title, req.Book.Author, req.Book.CategoryId, req.Book.Description, req.Book.PublishedDate, updatedAt).Scan(&updatedBook.Id, &updatedBook.Title, &updatedBook.Author, &updatedBook.CategoryId, &updatedBook.Description, &updatedBook.PublishedDate)
	if err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}

	return &api.UpdateBookResponse{
		Book: &updatedBook,
	}, nil
}

func (q *bookQuery) DeleteBook(ctx context.Context, tx pgx.Tx, req *api.DeleteBookRequest) (*api.DeleteBookResponse, error) {
	if req.BookId == "" {
		return nil, errors.New("book ID cannot be empty")
	}

	query := `UPDATE books SET deleted_at = $2 WHERE id = $1 AND deleted_at IS NULL`

	_, err := tx.Exec(ctx, query, req.BookId, time.Now())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("book with ID %s not found", req.BookId)
		}
		return nil, fmt.Errorf("failed to delete book: %w", err)
	}

	return &api.DeleteBookResponse{
		Success: true,
	}, nil
}
