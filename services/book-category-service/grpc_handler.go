package main

import (
	"context"
	"log"

	api "github.com/daffaromero/gobook/protobuf/api"
	"github.com/daffaromero/gobook/services/book-category-service/service"
	"google.golang.org/grpc"
)

type CategoryGRPCHandler struct {
	api.UnimplementedBookCategoryServiceServer

	service service.CategoryService
}

func NewCategoryGRPCHandler(server *grpc.Server, service service.CategoryService) {
	handler := &CategoryGRPCHandler{
		service: service,
	}

	api.RegisterBookCategoryServiceServer(server, handler)
}

func (h *CategoryGRPCHandler) GetCategory(ctx context.Context, req *api.GetCategoryRequest) (*api.GetCategoryResponse, error) {
	log.Printf("Received gRPC request for GetCategory: %v", req)
	res, err := h.service.GetCategory(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *CategoryGRPCHandler) ListCategories(ctx context.Context, req *api.ListCategoriesRequest) (*api.ListCategoriesResponse, error) {
	res, err := h.service.ListCategories(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
