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
	InitSQL  string `json:"init_sql"`
}

// response
type Data struct {
	Result          string
	Analysis        string
	Metrics         string
	Recommendations []string
	LLMAnswer       string
}

func (h *Handler) AnalyzeData(req AnalyzeRequest) (Data, error) {
	if req.SQLQuery == "" {
		return Data{}, errors.New("SQL Query is required")
	}

	// executing SQL with metrics
	result, metrics, err := h.svc.ExecuteSQLWithMetrics(req.SQLQuery, req.InitSQL)
	if err != nil {
		return Data{}, err
	}

	// request analysis
	analysis, err := h.svc.AnalyzeQueryInContainer(req.SQLQuery, req.InitSQL)
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
		LLMAnswer:       "",
	}

	model := "codellama12345"
	// model := "deepseek-r1:1.5b"
	llmanswer, err := services.QueryOllama(GeneratePrompt(d, req.InitSQL, req.SQLQuery), model)
	if err != nil {
		fmt.Println("OMG!! Error in QueryOllama")
	}

	return Data{
		Result:          result,
		Analysis:        analysis,
		Metrics:         metrics,
		Recommendations: recommendations,
		LLMAnswer:       llmanswer,
	}, nil
}

func (h *Handler) AnalyzeHandlerTemplate(c echo.Context) error {
	// reading request
	req := new(AnalyzeRequest)
	if err := c.Bind(req); err != nil {
		return c.String(400, "Invalid request body")
	}

	// executing analyzis
	data, err := h.AnalyzeData(*req)
	if err != nil {
		log.Printf("Error: %v", err)
		return c.JSON(400, echo.Map{"ExecutionError": err.Error()})
	}

	// returning result in html template
	return c.Render(http.StatusOK, "sql_analytics", data)
}

func (h *Handler) AnalyzeHandlerAPI(c echo.Context) error {
	// reading request
	req := new(AnalyzeRequest)
	if err := c.Bind(req); err != nil {
		return c.String(400, "Invalid request body")
	}

	// executing analyzis
	data, err := h.AnalyzeData(*req)
	if err != nil {
		log.Printf("Error: %v", err)
		return c.JSON(400, echo.Map{"ExecutionError": err.Error()})
	}

	// returning result in html template
	return c.JSON(http.StatusOK, data)
}
