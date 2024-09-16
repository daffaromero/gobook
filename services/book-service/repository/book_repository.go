package repository

import (
	"context"
	"fmt"

	api "github.com/daffaromero/gobook/protobuf/api"
	"github.com/daffaromero/gobook/services/book-service/repository/query"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepository interface {
	GetBook(ctx context.Context, req *api.GetBookRequest) (*api.GetBookResponse, error)
	ListBooks(ctx context.Context, req *api.ListBooksRequest) (*api.ListBooksResponse, error)
	CreateBook(ctx context.Context, req *api.CreateBookRequest) (*api.CreateBookResponse, error)
	UpdateBook(ctx context.Context, req *api.UpdateBookRequest) (*api.UpdateBookResponse, error)
	DeleteBook(ctx context.Context, id *api.DeleteBookRequest) (*api.DeleteBookResponse, error)
}

type bookRepository struct {
	db        Store
	bookQuery query.BookQuery
}

func NewBookRepository(db Store, bookQuery query.BookQuery) BookRepository {
	return &bookRepository{
		db:        db,
		bookQuery: bookQuery,
	}
}

func (r *bookRepository) GetBook(ctx context.Context, req *api.GetBookRequest) (*api.GetBookResponse, error) {
	var book *api.GetBookResponse

	err := r.db.WithoutTx(ctx, func(pool *pgxpool.Pool) error {
		var err error
		book, err = r.bookQuery.GetBook(ctx, req)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}
	return book, nil
}

func (r *bookRepository) ListBooks(ctx context.Context, req *api.ListBooksRequest) (*api.ListBooksResponse, error) {
	var books *api.ListBooksResponse

	err := r.db.WithoutTx(ctx, func(pool *pgxpool.Pool) error {
		var err error
		books, err = r.bookQuery.ListBooks(ctx, req)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list books: %w", err)
	}
	return books, nil
}

func (r *bookRepository) CreateBook(ctx context.Context, req *api.CreateBookRequest) (*api.CreateBookResponse, error) {
	var book *api.CreateBookResponse

	err := r.db.WithTx(ctx, func(tx pgx.Tx) error {
		var err error
		book, err = r.bookQuery.CreateBook(ctx, tx, req)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
	}
	return book, nil
}

func (r *bookRepository) UpdateBook(ctx context.Context, req *api.UpdateBookRequest) (*api.UpdateBookResponse, error) {
	var book *api.UpdateBookResponse

	err := r.db.WithTx(ctx, func(tx pgx.Tx) error {
		var err error
		book, err = r.bookQuery.UpdateBook(ctx, tx, req)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}
	return book, nil
}

func (r *bookRepository) DeleteBook(ctx context.Context, id *api.DeleteBookRequest) (*api.DeleteBookResponse, error) {
	var res *api.DeleteBookResponse

	err := r.db.WithTx(ctx, func(tx pgx.Tx) error {
		var err error
		res, err = r.bookQuery.DeleteBook(ctx, tx, id)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to delete book: %w", err)
	}
	return res, nil
}
