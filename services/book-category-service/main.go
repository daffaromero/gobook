package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/daffaromero/gobook/services/book-category-service/config"
	"github.com/daffaromero/gobook/services/book-category-service/controller"
	"github.com/daffaromero/gobook/services/book-category-service/repository"
	"github.com/daffaromero/gobook/services/book-category-service/repository/query"
	"github.com/daffaromero/gobook/services/book-category-service/service"
	"github.com/daffaromero/gobook/services/common/helper/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

var logs = logger.New("main")

func webServer() error {
	app := fiber.New()
	app.Use(requestid.New())

	serverConfig := config.NewServerConfig()
	dbConfig := config.NewPostgresDatabase()
	store := repository.NewStore(dbConfig)
	validate := validator.New()

	categoryQuery := query.NewCategoryQuery(dbConfig)
	categoryRepo := repository.NewCategoryRepository(store, categoryQuery)
	categoryService := service.NewCategoryService(categoryRepo, logs)
	categoryController := controller.NewCategoryController(validate, categoryService)

	app.Use(cors.New())

	categoryController.Route(app)

	err := app.Listen(serverConfig.Host, fiber.ListenConfig{
		DisableStartupMessage: true,
	})
	if err != nil {
		logs.Error("Failed to start server")
		return err
	}
	return nil
}

func main() {
	if err := webServer(); err != nil {
		logs.Error(err)
	}

	logs.Log("Web server started")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
}
