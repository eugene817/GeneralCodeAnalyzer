package handlers

import (
  "github.com/gofiber/fiber/v2"
  "github.com/eugene817/GeneralCodeAnalyzer/services"
)

// incoming request
type AnalyzeRequest struct {
  SQLQuerry string `json:"sql_querry"`
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


  result, err := services.ExecuteSQLInContainer(req.SQLQuerry, req.InitSQL)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(AnalyzeResponse{
      ExecutionErrors: []string{err.Error()},
    })
  }



  return c.JSON(AnalyzeResponse{
    Result: result,
  })
}


