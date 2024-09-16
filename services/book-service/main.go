package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/daffaromero/gobook/services/book-service/config"
	"github.com/daffaromero/gobook/services/book-service/controller"
	"github.com/daffaromero/gobook/services/book-service/repository"
	"github.com/daffaromero/gobook/services/book-service/repository/query"
	"github.com/daffaromero/gobook/services/book-service/service"
	"github.com/daffaromero/gobook/services/common/discovery"
	"github.com/daffaromero/gobook/services/common/discovery/consul"
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

	registry, err := consul.NewRegistry(serverConfig.ConsulAddr, serverConfig.Name)
	if err != nil {
		logs.Error("Failed to create consul registry for book service")
		return err
	}

	HTTPserviceID := discovery.GenerateServiceID(serverConfig.Name + "-http")
	httpPortInt, _ := strconv.Atoi(serverConfig.HTTPPort)

	ctx := context.Background()

	err = registry.RegisterService(ctx, serverConfig.Name, HTTPserviceID, serverConfig.HTTP, httpPortInt, []string{"http"})
	if err != nil {
		logs.Error("Failed to register HTTP book service to consul")
		return err
	}

	go func() {
		failureCount := 0
		const maxFailures = 5
		for {
			err := registry.HealthCheck(HTTPserviceID, serverConfig.Name)
			if err != nil {
				logs.Error(fmt.Sprintf("Failed to perform health check: %v", err))
				failureCount++
				if failureCount >= maxFailures {
					logs.Error("Max health check failures reached. Exiting health check loop.")
					break
				}
			} else {
				failureCount = 0
			}
			time.Sleep(time.Second * 2)
		}
	}()
	defer registry.DeregisterService(ctx, HTTPserviceID)

	bookQuery := query.NewBookQuery(dbConfig)
	bookRepo := repository.NewBookRepository(store, bookQuery)
	bookService, err := service.NewBookService(ctx, registry, bookRepo, logs)
	if err != nil {
		logs.Error("Failed to create book service")
		return err
	}
	bookController := controller.NewBookController(validate, bookService)

	logs.Log(fmt.Sprintf("Starting HTTP category server on %s", serverConfig.HTTP))
	app.Use(cors.New())
	bookController.Route(app)

	err = app.Listen(serverConfig.HTTP, fiber.ListenConfig{
		DisableStartupMessage: true,
	})
	if err != nil {
		logs.Error("Failed to start book server")
		return err
	}
	return nil
}

func main() {
	if err := webServer(); err != nil {
		logs.Error(err)
	}

	logs.Log("Book server started")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigchan
	logs.Log(fmt.Sprintf("Received signal: %s. Shutting down gracefully...", sig))
}
