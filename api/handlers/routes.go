package handlers

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {

	e.GET("/health", HealthCheckHandler)

	e.POST("/analyze", AnalyzeHandlerTemplate)

	e.POST("/analyze/json", AnalyzeHandlerAPI)

	e.POST("/analyze/python", AnalyzeHandlerTemplatePython)
}
