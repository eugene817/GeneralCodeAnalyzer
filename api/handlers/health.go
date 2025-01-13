package handlers

import (
	"github.com/labstack/echo/v4"
)

func HealthCheckHandler(c echo.Context) error {
  return c.JSON(404, echo.Map{
    "status": "ok",
  })
}
