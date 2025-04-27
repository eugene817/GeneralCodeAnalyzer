package handlers

import (
  "os"
	"github.com/labstack/echo/v4"
  "github.com/golang-jwt/jwt/v5"
  echojwt "github.com/labstack/echo-jwt/v4"
)

func (h *Handler) RegisterRoutes(e *echo.Echo) {

	e.GET("/health", HealthCheckHandler)

  e.POST("/register", h.Register)

  e.POST("/login", h.Login)

  // group up the routes for the JWT middleware
  r := e.Group("/api")
  configJWT := echojwt.Config{
    SigningKey: []byte(os.Getenv("JWT_SECRET")),
    TokenLookup: "cookie:jwt",
    ContextKey:  "user",
  }
  r.Use(echojwt.WithConfig(configJWT))

  r.GET("/private", func(c echo.Context) error {
    user := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)
    return c.JSON(200, echo.Map{"hello": user})
  })
  

	r.POST("/analyze/sql", h.AnalyzeHandlerTemplate)

	r.POST("/analyze/json", h.AnalyzeHandlerAPI)

	r.POST("/analyze/python", h.AnalyzeHandlerTemplatePython)

  r.POST("/analyze/c", h.AnalyzeHandlerTemplateC)



}
