package query

import (
	"context"

	api "github.com/daffaromero/gobook/protobuf/api"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryQuery interface {
	GetCategory(ctx context.Context, id *api.GetCategoryRequest) (*api.GetCategoryResponse, error)
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
