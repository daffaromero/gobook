package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	api "github.com/daffaromero/gobook/protobuf/api"
	"github.com/daffaromero/gobook/services/book-service/repository"
	"github.com/daffaromero/gobook/services/common/discovery"
	"github.com/daffaromero/gobook/services/common/helper/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BookService interface {
	GetBook(ctx context.Context, req *api.GetBookRequest) (*api.GetBookResponse, error)
	ListBooks(ctx context.Context, req *api.ListBooksRequest) (*api.ListBooksResponse, error)
	CreateBook(ctx context.Context, req *api.CreateBookRequest, title string, author string, categoryId string, description string) (*api.CreateBookResponse, error)
	UpdateBook(ctx context.Context, req *api.UpdateBookRequest, title string, author string, categoryId string, description string) (*api.UpdateBookResponse, error)
	DeleteBook(ctx context.Context, req *api.DeleteBookRequest) (*api.DeleteBookResponse, error)
}

type bookService struct {
	client   api.BookCategoryServiceClient
	registry discovery.Registry
	repo     repository.BookRepository
	logger   *logger.Log
}

func NewBookService(ctx context.Context, registry discovery.Registry, repo repository.BookRepository, logger *logger.Log) (*bookService, error) {
	conn, err := discovery.ServiceConnection(ctx, "book-category-service-grpc", registry)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to book-category-service: %v", err))

		return nil, err
	}
	logger.Log(fmt.Sprintf("Connected to book-category-service at %s", conn.Target()))

	client := api.NewBookCategoryServiceClient(conn)

	return &bookService{
		client:   client,
		registry: registry,
		repo:     repo,
		logger:   logger,
	}, nil
}

func (s *bookService) GetBook(ctx context.Context, req *api.GetBookRequest) (*api.GetBookResponse, error) {
	book, err := s.repo.GetBook(ctx, req)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to get book: %v", err))
		return nil, err
	}
	return book, nil
}

func (s *bookService) ListBooks(ctx context.Context, req *api.ListBooksRequest) (*api.ListBooksResponse, error) {
	books, err := s.repo.ListBooks(ctx, req)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to list books: %v", err))
		return nil, err
	}
	return books, nil
}

func (s *bookService) CreateBook(ctx context.Context, req *api.CreateBookRequest, title string, author string, categoryId string, description string) (*api.CreateBookResponse, error) {
	if s.client == nil {
		s.logger.Error("gRPC client is not initialized")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Book Category Service is not available.")
	}
	categoryReq := &api.GetCategoryRequest{
		CategoryId: categoryId,
	}
	category, err := s.client.GetCategory(ctx, categoryReq)
	if err != nil || category == nil || category.Category == nil {
		s.logger.Error(fmt.Sprintf("(RPC) Failed to get category: %v", err))
		return nil, err
	}

	now := &timestamppb.Timestamp{
		Seconds: time.Now().Unix(),
		Nanos:   int32(time.Now().Nanosecond()),
	}
	req.Book.Id = uuid.New().String()
	req.Book.Title = title
	req.Book.Author = author
	req.Book.CategoryId = category.Category.Id
	req.Book.Description = description
	req.Book.CreatedAt = now
	req.Book.UpdatedAt = now

	res, err := s.repo.CreateBook(ctx, req)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to create book: %v", err))
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Book already exists.")
		}
		return nil, err
	}

	return res, nil
}

func (s *bookService) UpdateBook(ctx context.Context, req *api.UpdateBookRequest, title string, author string, categoryId string, description string) (*api.UpdateBookResponse, error) {
	categoryReq := &api.GetCategoryRequest{
		CategoryId: categoryId,
	}

	category, err := s.client.GetCategory(ctx, categoryReq)
	if err != nil || category == nil || category.Category == nil {
		s.logger.Error(fmt.Sprintf("Failed to get category: %v", err))
		return nil, err
	}

	now := &timestamppb.Timestamp{
		Seconds: time.Now().Unix(),
		Nanos:   int32(time.Now().Nanosecond()),
	}
	req.Book.Title = title
	req.Book.Author = author
	req.Book.CategoryId = category.Category.Id
	req.Book.Description = description
	req.Book.UpdatedAt = now

	res, err := s.repo.UpdateBook(ctx, req)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to update book: %v", err))
		return nil, err
	}

	return res, nil
}

func (s *bookService) DeleteBook(ctx context.Context, req *api.DeleteBookRequest) (*api.DeleteBookResponse, error) {
	res, err := s.repo.DeleteBook(ctx, req)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to delete book: %v", err))
		return nil, err
	}
	return res, nil
}
