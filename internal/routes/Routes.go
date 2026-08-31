package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/kauanpecanha/odsquiz-quiz/internal/handlers"
	"github.com/kauanpecanha/odsquiz-quiz/internal/middleware"
	"github.com/kauanpecanha/odsquiz-quiz/internal/repositories"
	"github.com/kauanpecanha/odsquiz-quiz/internal/services"
	"github.com/kauanpecanha/odsquiz-quiz/pkg/config"
	"github.com/kauanpecanha/odsquiz-quiz/pkg/database"
)

func Setup(app *fiber.App, cfg *config.Config) {
	repository := repositories.NewRepo(database.DB.Db)
	service := &services.Service{Repo: repository}
	handler := &handlers.Handler{Service: service}

	app.Get("/", func(c fiber.Ctx) error { return c.SendString("ODS Quiz Quiz Microservice") })
	app.Get("/health", func(c fiber.Ctx) error { return c.SendString("ok") })
	app.Get("/api/quiz/questions", handler.ListQuestions)
	app.Post("/api/quiz/submissions", middleware.Protected(cfg.JWTSecret), handler.Submit)

	app.Use(func(c fiber.Ctx) error { return c.SendStatus(fiber.StatusNotFound) })
}
