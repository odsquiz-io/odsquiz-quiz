package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	"github.com/kauanpecanha/odsquiz-quiz/internal/migrations"
	"github.com/kauanpecanha/odsquiz-quiz/internal/routes"
	"github.com/kauanpecanha/odsquiz-quiz/pkg/config"
	"github.com/kauanpecanha/odsquiz-quiz/pkg/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	startupCtx, cancel := context.WithTimeout(context.Background(), database.StartupConnectionTimeout)
	defer cancel()

	err = migrations.RunMigrations(startupCtx, cfg)
	if err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName: "ODS Quiz Quiz Microservice",
	})

	app.Use(cors.New())

	routes.Setup(app, cfg)

	log.Fatal(app.Listen(":" + cfg.Port))
}
