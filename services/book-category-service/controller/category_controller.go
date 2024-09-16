package controller

import (
	api "github.com/daffaromero/gobook/protobuf/api"
	"github.com/daffaromero/gobook/services/book-category-service/config"
	"github.com/daffaromero/gobook/services/book-category-service/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type CategoryController interface {
	Route(*fiber.App)
	GetCategory(ctx fiber.Ctx) error
	ListCategories(ctx fiber.Ctx) error
	CreateCategory(ctx fiber.Ctx) error
	UpdateCategory(ctx fiber.Ctx) error
	DeleteCategory(ctx fiber.Ctx) error
}

type categoryController struct {
	validate *validator.Validate
	service  service.CategoryService
}

func NewCategoryController(validate *validator.Validate, service service.CategoryService) CategoryController {
	return &categoryController{
		validate: validate,
		service:  service,
	}
}

func (c *categoryController) Route(app *fiber.App) {
	api := app.Group(config.EndpointPrefix)
	api.Get("/:id", c.GetCategory)
	api.Get("/", c.ListCategories)
	api.Post("/new", c.CreateCategory)
	api.Put("/:id", c.UpdateCategory)
	api.Delete("/:id", c.DeleteCategory)
}

func (c *categoryController) GetCategory(ctx fiber.Ctx) error {
	var req api.GetCategoryRequest
	req.CategoryId = ctx.Params("id")
	if req.CategoryId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "category_id not provided"})
	}

	res, err := c.service.GetCategory(ctx.Context(), &req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *categoryController) ListCategories(ctx fiber.Ctx) error {
	var req *api.ListCategoriesRequest

	res, err := c.service.ListCategories(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *categoryController) CreateCategory(ctx fiber.Ctx) error {
	var req api.CreateCategoryRequest
	err := ctx.Bind().Body(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.validate.Struct(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.Category == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "category is required"})
	}

	name := req.Category.Name
	description := req.Category.Description

	res, err := c.service.CreateCategory(ctx.Context(), &req, name, description)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(res)
}

func (c *categoryController) UpdateCategory(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "category_id not provided"})
	}

	var req api.UpdateCategoryRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.Category == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "category cannot be empty"})
	}
	req.Category.Id = id

	if err := c.validate.Struct(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.Category.Name == "" && req.Category.Description == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "At least one field (name or description) must be provided for update"})
	}

	res, err := c.service.UpdateCategory(ctx.Context(), &req, req.Category.Name, req.Category.Description)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *categoryController) DeleteCategory(ctx fiber.Ctx) error {
	var req api.DeleteCategoryRequest
	req.CategoryId = ctx.Params("id")
	if req.CategoryId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "category_id not provided"})
	}

	res, err := c.service.DeleteCategory(ctx.Context(), &req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}
