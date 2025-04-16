package handlers

import (
	"encoding/json"
)

type AnalysisRequest struct {
	Result          string            `json:"result"`
	Analysis        string            `json:"analysis"`
	Metrics         map[string]string `json:"metrics"`
	Recommendations []string          `json:"recommendations"`
}

func toJson(v interface{}) (string, error) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func GeneratePrompt(req Data, sqlQuerry, initSQL string) string {
	metricsJson, _ := json.Marshal(req.Metrics)
	recommendationsJson, _ := json.Marshal(req.Recommendations)

	return `
  Analyze the following SQL query and its result, give complicated answers with more recommendations:

  SQL Query:
  ` + sqlQuerry + `

  Initial SQL:
  ` + initSQL + `

  Result:
  ` + req.Result + `

  Analysis:
  ` + req.Analysis + `

  Metrics:
  ` + string(metricsJson) + `

  Recommendations:
  ` + string(recommendationsJson) + `
  `
}

func GeneratePromptPython(req PythonData, pythonCode string) string {
	metricsJson, _ := json.Marshal(req.Metrics)
	recommendationsJson, _ := json.Marshal(req.Recommendations)

	return `
  Analyze the following python code giving recommendations

  Python Code:
  ` + pythonCode + `

  Result:
  ` + req.Result + `

  Metrics:
  ` + string(metricsJson) + `

  Recommendations:
  ` + string(recommendationsJson) + `
  `
}
