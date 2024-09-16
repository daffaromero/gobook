package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	api "github.com/daffaromero/gobook/protobuf/api"
	"github.com/daffaromero/gobook/services/book-category-service/repository"
	"github.com/daffaromero/gobook/services/common/helper/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CategoryService interface {
	GetCategory(ctx context.Context, req *api.GetCategoryRequest) (*api.GetCategoryResponse, error)
	ListCategories(ctx context.Context, req *api.ListCategoriesRequest) (*api.ListCategoriesResponse, error)
	CreateCategory(ctx context.Context, req *api.CreateCategoryRequest, name string, desc string) (*api.CreateCategoryResponse, error)
	UpdateCategory(ctx context.Context, req *api.UpdateCategoryRequest, name string, desc string) (*api.UpdateCategoryResponse, error)
	DeleteCategory(ctx context.Context, req *api.DeleteCategoryRequest) (*api.DeleteCategoryResponse, error)
}

type categoryService struct {
	repo   repository.CategoryRepository
	logger *logger.Log
}

func NewCategoryService(repo repository.CategoryRepository, logger *logger.Log) CategoryService {
	return &categoryService{
		repo:   repo,
		logger: logger,
	}
}

func (s *categoryService) GetCategory(ctx context.Context, req *api.GetCategoryRequest) (*api.GetCategoryResponse, error) {
	category, err := s.repo.GetCategory(ctx, req)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to get category: %v", err))
		return nil, err
	}
	return category, nil
}

func (s *categoryService) ListCategories(ctx context.Context, req *api.ListCategoriesRequest) (*api.ListCategoriesResponse, error) {
	categories, err := s.repo.ListCategories(ctx, req)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to list categories: %v", err))
		return nil, err
	}
	return categories, nil
}

func (s *categoryService) CreateCategory(ctx context.Context, req *api.CreateCategoryRequest, name string, desc string) (*api.CreateCategoryResponse, error) {

	now := &timestamppb.Timestamp{
		Seconds: time.Now().Unix(),
		Nanos:   int32(time.Now().Nanosecond()),
	}
	req.Category.Id = uuid.New().String()
	req.Category.Name = name
	req.Category.Description = desc
	req.Category.CreatedAt = now
	req.Category.UpdatedAt = now

	res, err := s.repo.CreateCategory(ctx, req)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to create category: %v", err))
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Category already exists.")
		}
		return nil, err
	}

	return res, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, req *api.UpdateCategoryRequest, name string, desc string) (*api.UpdateCategoryResponse, error) {
	now := &timestamppb.Timestamp{
		Seconds: time.Now().Unix(),
		Nanos:   int32(time.Now().Nanosecond()),
	}
	req.Category.Name = name
	req.Category.Description = desc
	req.Category.UpdatedAt = now

	res, err := s.repo.UpdateCategory(ctx, req)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to update category: %v", err))
		return nil, err
	}

	return res, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, req *api.DeleteCategoryRequest) (*api.DeleteCategoryResponse, error) {
	res, err := s.repo.DeleteCategory(ctx, req)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to delete category: %v", err))
		return nil, err
	}

	return res, nil
}
