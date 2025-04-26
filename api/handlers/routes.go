package handlers

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) RegisterRoutes(e *echo.Echo) {

	e.GET("/health", HealthCheckHandler)

	e.POST("/analyze", h.AnalyzeHandlerTemplate)

	e.POST("/analyze/json", h.AnalyzeHandlerAPI)

	e.POST("/analyze/python", h.AnalyzeHandlerTemplatePython)
}
