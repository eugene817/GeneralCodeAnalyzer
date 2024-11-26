package handlers

import (
  "github.com/gofiber/fiber/v2"
  "github.com/eugene817/GeneralCodeAnalyzer/services"
  "log"
)

// incoming request
type AnalyzeRequest struct {
  SQLQuery string `json:"sql_query"`
  InitSQL string `json:"init_sql"`
}

// response
type AnalyzeResponse struct {
  Result string `json:"result"`
  ExecutionErrors []string `json:"execution_errors, omitempty"`
}


func AnalyzeHandler(c *fiber.Ctx) error {
  req := new(AnalyzeRequest)
  
  if err := c.BodyParser(req); err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error": "Invalid request body",
    })
  }
  log.Printf("SQL Query: %s", req.SQLQuery)
  log.Printf("Init SQL: %s", req.InitSQL)

  if req.SQLQuery == "" {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error": "SQL Query is required",
    })
  }

  result, metrics, err := services.ExecuteSQLWithMemoryMetrics(req.SQLQuery, req.InitSQL)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(AnalyzeResponse{
      ExecutionErrors: []string{err.Error()},
    })
  }

  analysis, err := services.AnalyzeQueryInContainer(req.SQLQuery, req.InitSQL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

  recommendations := services.GenerateRecommendations(metrics, analysis, req.SQLQuery)

  return c.JSON(fiber.Map{
    "result": result,
    "analysis": analysis,
    "metrics": metrics,
    "recommendations": recommendations,
  })
}


