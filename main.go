package main

import (
	"github.com/eugene817/GeneralCodeAnalyzer/api/handlers"
	"github.com/eugene817/GeneralCodeAnalyzer/api/templates"
	"github.com/eugene817/GeneralCodeAnalyzer/config"
	"github.com/eugene817/Cownets"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	handlers.RegisterRoutes(e)
  templates.RegisterTemplatesRoutes(e)

	e.Static("/static", "./api/static")
	port := config.GetPort()

	e.Logger.Fatal(e.Start(port))
}
