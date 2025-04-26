package main

import (
  "log"
	"fmt"
	"os"
  "time"

	"github.com/eugene817/Cowdocs/api"
	"github.com/eugene817/Cowdocs/container"
	"github.com/eugene817/GeneralCodeAnalyzer/api/handlers"
	"github.com/eugene817/GeneralCodeAnalyzer/api/templates"
	"github.com/eugene817/GeneralCodeAnalyzer/config"
	"github.com/eugene817/GeneralCodeAnalyzer/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

  mgr, err := container.NewDockerManager()
	if err != nil {
		fmt.Errorf("failed to create Docker manager: %v", err)
    os.Exit(1)
	}

	mng := api.NewAPI(mgr)

  // Ensure the images are available
  Images := []string{
    "python:3",
    "keinos/sqlite3",
  }
  
  for {
        if err := mng.Ping(); err == nil {
            break
        }
        log.Println("waiting for Docker daemon…")
        time.Sleep(500 * time.Millisecond)
  }

  if err := mng.EnsureImages(Images); err != nil {
    log.Fatalf("failed to pull initial images: %v", err)    
    os.Exit(1)
  }

	e := echo.New()

	e.Use(middleware.Logger())

  svc := services.NewService(mng)
  h := handlers.NewHandler(svc)
  h.RegisterRoutes(e)

	templates.RegisterTemplatesRoutes(e)

	e.Static("/static", "./api/static")
	port := config.GetPort()

	e.Logger.Fatal(e.Start(port))
}
