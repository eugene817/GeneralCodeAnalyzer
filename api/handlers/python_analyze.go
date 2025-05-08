package handlers

import (
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/eugene817/GeneralCodeAnalyzer/services"
	"github.com/labstack/echo/v4"
)

// incoming request
type PythonAnalyzeRequest struct {
	PythonCode string `json:"python_code"`
}

// response
type PythonData struct {
	Result          string
	Metrics         string
	Recommendations []string
	LLMAnswer       string
}

func (h *Handler) PythonAnalyzeData(req PythonAnalyzeRequest) (PythonData, error) {
	if req.PythonCode == "" {
		return PythonData{}, errors.New("Python code is empty")
	}

	// executing SQL with metrics
	result, metrics, err := h.svc.ExecutePythonWithMetrics(req.PythonCode)
	if err != nil {
		return PythonData{}, err
	}

	// SQL (not used yet)
	_, err = h.svc.ExecutePythonInContainer(req.PythonCode)
	if err != nil {
		return PythonData{}, err
	}

	// recomendations generation
	recommendations := services.GenerateRecommendationsPython(metrics, req.PythonCode)

	d := PythonData{
		Result:          result,
		Metrics:         metrics,
		Recommendations: recommendations,
		LLMAnswer:       "",
	}

	model := "codellama1234"
	// model := "deepseek-r1:1.5b"
	llmanswer, err := services.QueryOllama(GeneratePromptPython(d, req.PythonCode), model)
	if err != nil {
		fmt.Println("OMG!! Error in QueryOllama")
	}

	return PythonData{
		Result:          result,
		Metrics:         metrics,
		Recommendations: recommendations,
		LLMAnswer:       llmanswer,
	}, nil
}

func (h *Handler) AnalyzeHandlerTemplatePython(c echo.Context) error {
	// reading request
	req := new(PythonAnalyzeRequest)
	if err := c.Bind(req); err != nil {
		return c.String(400, "Invalid request body")
	}

	// executing analyzis
	data, err := h.PythonAnalyzeData(*req)
	if err != nil {
		log.Printf("Error: %v", err)
		return c.JSON(400, echo.Map{"ExecutionError": err.Error()})
	}

	// returning result in html template
	return c.Render(http.StatusOK, "python_analytics", data)
}

// --------- Lint ------------

type PythonLintRequest struct {
	PythonCode string `json:"python_code"`
}

func (h *Handler) PythonLintHandler(c echo.Context) error {
	req := new(PythonLintRequest)
	if err := c.Bind(req); err != nil {
		return c.HTML(http.StatusOK,
			`<pre class="text-red-600">Invalid payload</pre>`)
	}

	diag, err := h.svc.LintPythonInContainer(req.PythonCode)
	if err != nil {
		return c.HTML(http.StatusOK,
			`<pre class="text-red-600">`+
				html.EscapeString(err.Error())+`</pre>`)
	}

	if diag == "" {
		diag = "No syntax errors."
	}
	return c.HTML(http.StatusOK,
		`<pre class="text-green-600">`+
			html.EscapeString(diag)+`</pre>`)
}
