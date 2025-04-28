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
type CAnalyzeRequest struct {
	CCode string `json:"c_code"`
}

// response
type CData struct {
	Result          string
	Metrics         string
	Recommendations []string
	LLMAnswer       string
}

func (h *Handler) CAnalyzeData(req CAnalyzeRequest) (CData, error) {
	if req.CCode == "" {
		return CData{}, errors.New("C code is empty")
	}

	// executing C with metrics
	result, metrics, err := h.svc.ExecuteCWithMetrics(req.CCode)
	if err != nil {
		return CData{}, err
	}

	// SQL (not used yet)
	_, err = h.svc.ExecuteCInContainer(req.CCode)
	if err != nil {
		return CData{}, err
	}

	// recomendations generation
	recommendations := services.GenerateRecommendationsC(metrics, req.CCode)

	d := CData{
		Result:          result,
		Metrics:         metrics,
		Recommendations: recommendations,
		LLMAnswer:       "",
	}

	model := "codellama1234"
	// model := "deepseek-r1:1.5b"
	llmanswer, err := services.QueryOllama(GeneratePromptC(d, req.CCode), model)
	if err != nil {
		fmt.Println("OMG!! Error in QueryOllama")
	}

	return CData{
		Result:          result,
		Metrics:         metrics,
		Recommendations: recommendations,
		LLMAnswer:       llmanswer,
	}, nil
}

func (h *Handler) AnalyzeHandlerTemplateC(c echo.Context) error {
	// reading request
	req := new(CAnalyzeRequest)
	if err := c.Bind(req); err != nil {
		return c.String(400, "Invalid request body")
	}

	// executing analyzis
	data, err := h.CAnalyzeData(*req)
	if err != nil {
		log.Printf("Error: %v", err)
		return c.JSON(400, echo.Map{"ExecutionError": err.Error()})
	}
  log.Printf("data: %v", data)

	// returning result in html template
	return c.Render(http.StatusOK, "c_analytics", data)
}

// --------- Lint ------------

type CLintRequest struct {
	CCode string `json:"c_code"`
}

// CLintResponse
type CLintResponse struct {
	Diagnostics string `json:"diagnostics"`
}

func (h *Handler) CLintHandler(c echo.Context) error {
	req := new(CLintRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid payload"})
	}
	diag, err := h.svc.LintCInContainer(req.CCode)
	if err != nil {
		log.Printf("lint error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if diag == "" {
		diag = "No syntax errors."
	}
	return c.JSON(http.StatusOK, CLintResponse{Diagnostics: diag})
}

