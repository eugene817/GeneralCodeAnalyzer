package main

import (
	"net/http"

	"github.com/eugene817/GeneralCodeAnalyzer/api/handlers"
	"github.com/eugene817/GeneralCodeAnalyzer/api/templates"
	"github.com/eugene817/GeneralCodeAnalyzer/config"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	handlers.RegisterRoutes(e)

	e.Renderer = templates.NewTemplate()
	e.Static("/static", "./api/static")

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})

	port := config.GetPort()

	e.Logger.Fatal(e.Start(port))
}
