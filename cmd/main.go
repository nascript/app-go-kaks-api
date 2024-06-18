package main

import (
	"context"
	"errors"
	"kaks-cloud-web-api-task/api/product"
	// docs
	_ "kaks-cloud-web-api-task/docs"
	"kaks-cloud-web-api-task/infrastructure"
	storeRepo "kaks-cloud-web-api-task/repository/product"
	storeServ "kaks-cloud-web-api-task/service/product"

	otelfiber "github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	middlewareLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger" // Pastikan ini diimport
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"

	"log"
	"os"

	"golang.org/x/exp/slog"
)

// @title kaks-cloud-web-api-task
// @version 1.0
// @Author https://github.com/nascript
// @description This web-service-api kaks-cloud-web-api-task. this project for the KAKS master's degree in computer science 
// @termsOfService http://swagger.io/terms/
// @contact.name Nasution
// @contact.url https://www.github.io/nascript
// @contact.email nasutioncode@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 3.27.69.11:4040
// @BasePath /

func main() {
	// Initialize Viper
	viper.AutomaticEnv() // Read environment variables
	ctx := context.Background()

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return a
		},
	}).WithAttrs([]slog.Attr{
		slog.String("service", os.Getenv("OTEL_SERVICE_NAME")),
		slog.String("with-release", "v1.0.0"),
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	// Set up OpenTelemetry.
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	otelCollector := os.Getenv("OTEL_COLLECTOR")
	serviceVersion := "0.1.0"
	otelShutdown, err := infrastructure.SetupOTelSDK(ctx, otelCollector, serviceName, serviceVersion, os.Getenv("OTEL_ENV"))
	if err != nil {
		log.Fatalf("failed to initialize OTel SDK: %v", err)
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// init mongo
	mongo := infrastructure.NewMongo(ctx, os.Getenv("MONGO_DSN"), os.Getenv("MONGO_DB_NAME"))
	mongo = mongo.Connect()

	// store
	storeRepository := storeRepo.NewstoreRepository(mongo.Client, mongo.DB, "products")
	storeService := storeServ.NewStoreService(storeRepository)
	handler := product.NewStoreHandler(storeService)

	app := fiber.New()

	app.Use(middlewareLogger.New(middlewareLogger.Config{
		Format: "[${time}] ${ip}  ${status} - ${latency} ${method} ${path}\n",
	}))

	app.Use(otelfiber.Middleware())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	v1 := app.Group("/api/v1")
	v1.Get("/swagger/*", swagger.HandlerDefault)

	v1.Get("/product/:id", handler.Get)
	v1.Get("/product", handler.GetAll)
	v1.Post("/product", handler.Create)
	v1.Put("/product/:id", handler.Update)
	v1.Delete("/product/:id", handler.Delete)

	app.Listen(":4040")
}
