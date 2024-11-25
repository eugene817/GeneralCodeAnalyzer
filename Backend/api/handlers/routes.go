package handlers

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App) {
	app.Get("/health", HealthCheckHandler)

	app.Post("/analyze", AnalyzeHandler)
}
