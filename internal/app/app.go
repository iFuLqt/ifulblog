package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ifulblog/config"
	"ifulblog/lib/auth"
	"ifulblog/lib/middleware"
	"ifulblog/lib/pagination"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func RunServer() {
	cfg := config.NewConfig()
	_, err := cfg.ConnectionPostgres()
	if err != nil {
		log.Fatal("Error connecting to database %v", err)
		return
	}

	//cloudfareR2
	cdfR2 := cfg.LoadAwsConfig()
	_ = s3.NewFromConfig(cdfR2)

	_ = auth.NewJwt(cfg)
	_ = middleware.NewMiddleware(cfg)
	_ = pagination.NewPagination()

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] %{ip} %{status} - %{latency} %{method} %{path}\n",
	}))

	_ = app.Group("/api")

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppEnv = os.Getenv("APP_PORT")
		}
		err := app.Listen(":" + cfg.App.AppPort)
		if err != nil {
			log.Fatal("error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	<-quit

	log.Println("server shutdown of 5 seconds")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.ShutdownWithContext(ctx)
}
