package handlers

import "github.com/gofiber/fiber/v2"

func HealthCheckHandler(c *fiber.Ctx) error {
  return c.JSON(fiber.Map{
    "status": "ok",
  })
}
