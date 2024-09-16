package repository

import (
	"context"
	"fmt"

	api "github.com/daffaromero/gobook/protobuf/api"
	"github.com/daffaromero/gobook/services/book-category-service/repository/query"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository interface {
	GetCategory(ctx context.Context, req *api.GetCategoryRequest) (*api.GetCategoryResponse, error)
	ListCategories(ctx context.Context, req *api.ListCategoriesRequest) (*api.ListCategoriesResponse, error)
	CreateCategory(ctx context.Context, req *api.CreateCategoryRequest) (*api.CreateCategoryResponse, error)
	UpdateCategory(ctx context.Context, req *api.UpdateCategoryRequest) (*api.UpdateCategoryResponse, error)
	DeleteCategory(ctx context.Context, id *api.DeleteCategoryRequest) (*api.DeleteCategoryResponse, error)
}

type categoryRepository struct {
	db            Store
	categoryQuery query.CategoryQuery
}

func NewCategoryRepository(db Store, categoryQuery query.CategoryQuery) CategoryRepository {
	return &categoryRepository{
		db:            db,
		categoryQuery: categoryQuery,
	}
}

func (r *categoryRepository) GetCategory(ctx context.Context, req *api.GetCategoryRequest) (*api.GetCategoryResponse, error) {
	var category *api.GetCategoryResponse

	err := r.db.WithoutTx(ctx, func(pool *pgxpool.Pool) error {
		var err error
		category, err = r.categoryQuery.GetCategory(ctx, req)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return category, nil
}

func (r *categoryRepository) ListCategories(ctx context.Context, req *api.ListCategoriesRequest) (*api.ListCategoriesResponse, error) {
	var categories *api.ListCategoriesResponse

	err := r.db.WithoutTx(ctx, func(pool *pgxpool.Pool) error {
		var err error
		categories, err = r.categoryQuery.ListCategories(ctx, req)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	return categories, nil
}

func (r *categoryRepository) CreateCategory(ctx context.Context, req *api.CreateCategoryRequest) (*api.CreateCategoryResponse, error) {
	var category *api.CreateCategoryResponse

	err := r.db.WithTx(ctx, func(tx pgx.Tx) error {
		var err error
		category, err = r.categoryQuery.CreateCategory(ctx, tx, req)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	return category, nil
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, req *api.UpdateCategoryRequest) (*api.UpdateCategoryResponse, error) {
	var category *api.UpdateCategoryResponse

	err := r.db.WithTx(ctx, func(tx pgx.Tx) error {
		var err error
		category, err = r.categoryQuery.UpdateCategory(ctx, tx, req)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}
	return category, nil
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, req *api.DeleteCategoryRequest) (*api.DeleteCategoryResponse, error) {
	var res *api.DeleteCategoryResponse

	err := r.db.WithTx(ctx, func(tx pgx.Tx) error {
		var err error
		res, err = r.categoryQuery.DeleteCategory(ctx, tx, req)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to delete category: %w", err)
	}
	return res, nil
}
