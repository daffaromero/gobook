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

type CategoryQuery interface {
	GetCategory(ctx context.Context, req *api.GetCategoryRequest) (*api.GetCategoryResponse, error)
	ListCategories(ctx context.Context, req *api.ListCategoriesRequest) (*api.ListCategoriesResponse, error)
	CreateCategory(ctx context.Context, tx pgx.Tx, req *api.CreateCategoryRequest) (*api.CreateCategoryResponse, error)
	UpdateCategory(ctx context.Context, tx pgx.Tx, req *api.UpdateCategoryRequest) (*api.UpdateCategoryResponse, error)
	DeleteCategory(ctx context.Context, tx pgx.Tx, id *api.DeleteCategoryRequest) (*api.DeleteCategoryResponse, error)
}

type categoryQuery struct {
	db *pgxpool.Pool
}

func NewCategoryQuery(db *pgxpool.Pool) *categoryQuery {
	return &categoryQuery{
		db: db,
	}
}

func (q *categoryQuery) GetCategory(ctx context.Context, req *api.GetCategoryRequest) (*api.GetCategoryResponse, error) {
	if req == nil || req.CategoryId == "" {
		return nil, errors.New("category ID cannot be empty")
	}
	query := `SELECT id, name, description FROM book_categories WHERE id = $1`

	row := q.db.QueryRow(ctx, query, req.CategoryId)

	var category api.BookCategory
	err := row.Scan(&category.Id, &category.Name, &category.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("category with ID %s not found", req.CategoryId)
		}
		return nil, fmt.Errorf("failed to scan category: %w", err)
	}

	return &api.GetCategoryResponse{
		Category: &category,
	}, nil
}

func (q *categoryQuery) ListCategories(ctx context.Context, req *api.ListCategoriesRequest) (*api.ListCategoriesResponse, error) {
	query := `SELECT id, name, description FROM book_categories WHERE deleted_at IS NULL`

	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []*api.BookCategory
	for rows.Next() {
		var category api.BookCategory
		err := rows.Scan(&category.Id, &category.Name, &category.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, &category)
	}

	return &api.ListCategoriesResponse{
		Categories: categories,
	}, nil
}

func (q *categoryQuery) CreateCategory(ctx context.Context, tx pgx.Tx, req *api.CreateCategoryRequest) (*api.CreateCategoryResponse, error) {
	if req == nil || req.Category == nil {
		return nil, errors.New("request cannot be nil")
	}

	query := `INSERT INTO book_categories (id, name, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, name, description`

	createdAt := req.Category.CreatedAt.AsTime()
	updatedAt := req.Category.UpdatedAt.AsTime()

	err := tx.QueryRow(ctx, query, req.Category.Id, req.Category.Name, req.Category.Description, createdAt, updatedAt).Scan(&req.Category.Id, &req.Category.Name, &req.Category.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to insert category: %w", err)
	}

	return &api.CreateCategoryResponse{
		Category: &api.BookCategory{
			Id:          req.Category.Id,
			Name:        req.Category.Name,
			Description: req.Category.Description,
		},
	}, nil
}

func (q *categoryQuery) UpdateCategory(ctx context.Context, tx pgx.Tx, req *api.UpdateCategoryRequest) (*api.UpdateCategoryResponse, error) {
	if req == nil || req.Category == nil {
		return nil, errors.New("request and category cannot be nil")
	}
	if req.Category.Id == "" {
		return nil, errors.New("category ID cannot be empty")
	}

	query := `UPDATE book_categories 
		SET 
			name = COALESCE($2, name),
			description = COALESCE($3, description),
			updated_at = $4
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, name, description, updated_at`

	updatedAt := time.Now()
	if req.Category.UpdatedAt != nil {
		updatedAt = req.Category.UpdatedAt.AsTime()
	}

	var updatedCategory api.BookCategory

	err := tx.QueryRow(ctx, query, req.Category.Id, req.Category.Name, req.Category.Description, updatedAt).Scan(&updatedCategory.Id, &updatedCategory.Name, &updatedCategory.Description, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("category with ID %s not found", req.Category.Id)
		}
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return &api.UpdateCategoryResponse{
		Category: &updatedCategory,
	}, nil
}

func (q *categoryQuery) DeleteCategory(ctx context.Context, tx pgx.Tx, req *api.DeleteCategoryRequest) (*api.DeleteCategoryResponse, error) {
	if req.CategoryId == "" {
		return nil, errors.New("category ID cannot be empty")
	}

	query := `UPDATE book_categories SET deleted_at = $2 WHERE id = $1 AND deleted_at IS NULL`

	_, err := tx.Exec(ctx, query, req.CategoryId, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to delete category: %w", err)
	}

	return &api.DeleteCategoryResponse{
		Success: true,
	}, nil
}
