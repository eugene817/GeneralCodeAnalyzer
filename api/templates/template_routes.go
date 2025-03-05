package templates

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterTemplatesRoutes(e *echo.Echo) {
  e.Renderer = NewTemplate()

  e.GET("/", func(c echo.Context) error {
    return c.Render(http.StatusOK, "index", nil)
  })

  e.GET("/sql", func(c echo.Context) error {
    return c.Render(http.StatusOK, "sql", nil)
  })

  e.GET("/python", func(c echo.Context) error {
    return c.Render(http.StatusOK, "python", nil)
  })
}
