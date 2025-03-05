package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/eugene817/GeneralCodeAnalyzer/services"
	"github.com/labstack/echo/v4"
)

// incoming request
type AnalyzeRequest struct {
  SQLQuery string `json:"sql_query"`
  InitSQL string `json:"init_sql"`
}

// response
type Data struct {
	Result          string
	Analysis        string
	Metrics         map[string]interface{}
	Recommendations []string
  LLMAnswer       string
}


func AnalyzeData(req AnalyzeRequest) (Data, error) {
	if req.SQLQuery == "" {
		return Data{}, errors.New("SQL Query is required")
	}

	// executing SQL with metrics
	result, metrics, err := services.ExecuteSQLWithMemoryMetrics(req.SQLQuery, req.InitSQL)
	if err != nil {
		return Data{}, err
	}

	// SQL (not used yet)
	_, err = services.ExecuteSQLInContainer(req.SQLQuery, req.InitSQL)
	if err != nil {
    return Data{}, err
	}


	// request analysis
	analysis, err := services.AnalyzeQueryInContainer(req.SQLQuery, req.InitSQL)
	if err != nil {
		return Data{}, err
	}

	// recomendations generation
	recommendations := services.GenerateRecommendations(metrics, analysis, req.SQLQuery)

  d := Data{
		Result:          result,
		Analysis:        analysis,
		Metrics:         metrics,
		Recommendations: recommendations,
    LLMAnswer: "",
	}

  //model := "codellama"
  model := "deepseek-r1:1.5b"
  llmanswer, err := services.QueryOllama(GeneratePrompt(d, req.InitSQL, req.SQLQuery), model)
  if err != nil {
    fmt.Println("OMG!! Error in QueryOllama")    
  }


	return Data{
		Result:          result,
		Analysis:        analysis,
		Metrics:         metrics,
		Recommendations: recommendations,
    LLMAnswer: llmanswer,
	}, nil
}


func AnalyzeHandlerTemplate(c echo.Context) error {
	// reading request
	req := new(AnalyzeRequest)
	if err := c.Bind(req); err != nil {
		return c.String(400, "Invalid request body")
	}

	// executing analyzis
	data, err := AnalyzeData(*req)
	if err != nil {
		log.Printf("Error: %v", err)
		return c.JSON(400, echo.Map{"ExecutionError": err.Error()})
	}

  // returning result in html template
	return c.Render(http.StatusOK, "analytics", data)
}


func AnalyzeHandlerAPI(c echo.Context) error {
	// reading request
	req := new(AnalyzeRequest)
	if err := c.Bind(req); err != nil {
		return c.String(400, "Invalid request body")
	}


	// executing analyzis
	data, err := AnalyzeData(*req)
	if err != nil {
		log.Printf("Error: %v", err)
		return c.JSON(400, echo.Map{"ExecutionError": err.Error()})
	}

  // returning result in html template
	return c.JSON(http.StatusOK, data)
}

