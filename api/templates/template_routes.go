package templates

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterTemplatesRoutes(e *echo.Echo) {
  e.Renderer = NewTemplate()

  e.GET("/", func(c echo.Context) error {
    return c.Render(http.StatusOK, "base", echo.Map{
      "Title": "Home",
      "Content": "<h1>Welcome to the Home Page</h1>",
    })
  })

  e.GET("/login", func(c echo.Context) error {
    return c.Render(http.StatusOK, "login_partial", echo.Map{
      "Error": c.QueryParam("error"),
    })
  })
  e.GET("/register", func(c echo.Context) error {
    return c.Render(http.StatusOK, "register_partial", echo.Map{
      "Error": c.QueryParam("error"),
    })
  })

  r := e.Group("/api")


  r.GET("/sql", func(c echo.Context) error {
    return c.Render(http.StatusOK, "sql", nil)
  })

  r.GET("/python", func(c echo.Context) error {
    return c.Render(http.StatusOK, "python", nil)
  })

  r.GET("/c", func(c echo.Context) error {
    return c.Render(http.StatusOK, "c", nil)
  })
}
