package main

import (
	"log"

	"github.com/eugene817/GeneralCodeAnalyzer/api/handlers"
	"github.com/eugene817/GeneralCodeAnalyzer/config"
	"github.com/gofiber/fiber/v2"
)

func main() {
  app := fiber.New()

  handlers.RegisterRoutes(app)

  port := config.GetPort()
  log.Printf("Server running at http://localhost%s", &port)
  log.Fatal(app.Listen(port))
}
