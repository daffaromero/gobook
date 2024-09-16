package controller

import (
	api "github.com/daffaromero/gobook/protobuf/api"
	"github.com/daffaromero/gobook/services/book-service/config"
	"github.com/daffaromero/gobook/services/book-service/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type BookController interface {
	Route(*fiber.App)
	GetBook(ctx fiber.Ctx) error
	ListBooks(ctx fiber.Ctx) error
	CreateBook(ctx fiber.Ctx) error
	UpdateBook(ctx fiber.Ctx) error
	DeleteBook(ctx fiber.Ctx) error
}

type bookController struct {
	validate *validator.Validate
	service  service.BookService
}

func NewBookController(validate *validator.Validate, service service.BookService) BookController {
	return &bookController{
		validate: validate,
		service:  service,
	}
}

func (c *bookController) Route(app *fiber.App) {
	api := app.Group(config.EndpointPrefix)
	api.Get("/:id", c.GetBook)
	api.Get("/", c.ListBooks)
	api.Post("/new", c.CreateBook)
	api.Put("/:id", c.UpdateBook)
	api.Delete("/:id", c.DeleteBook)
}

func (c *bookController) GetBook(ctx fiber.Ctx) error {
	var req api.GetBookRequest
	req.BookId = ctx.Params("id")
	if req.BookId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "book_id not provided"})
	}

	res, err := c.service.GetBook(ctx.Context(), &req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *bookController) ListBooks(ctx fiber.Ctx) error {
	var req *api.ListBooksRequest

	res, err := c.service.ListBooks(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *bookController) CreateBook(ctx fiber.Ctx) error {
	var req api.CreateBookRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.validate.Struct(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.Book == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "book cannot be empty"})
	}

	res, err := c.service.CreateBook(ctx.Context(), &req, req.Book.Title, req.Book.Author, req.Book.CategoryId, req.Book.Description)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(res)
}

func (c *bookController) UpdateBook(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "book_id not provided"})
	}

	var req api.UpdateBookRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.validate.Struct(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.Book == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "book cannot be empty"})
	}
	req.Book.Id = id

	if err := c.validate.Struct(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.Book.Author == "" && req.Book.CategoryId == "" && req.Book.Description == "" && req.Book.Title == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no fields to update"})
	}

	res, err := c.service.UpdateBook(ctx.Context(), &req, req.Book.Title, req.Book.Author, req.Book.CategoryId, req.Book.Description)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *bookController) DeleteBook(ctx fiber.Ctx) error {
	var req api.DeleteBookRequest
	req.BookId = ctx.Params("id")
	if req.BookId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "book_id not provided"})
	}

	res, err := c.service.DeleteBook(ctx.Context(), &req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}
