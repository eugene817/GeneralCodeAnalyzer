package analyzer

import (
  "errors"
  "fmt"
  "github.com/eugene817/GeneralCodeAnalyzer/pkg/analyzer/languages/python/static"

)

// Analysis results
type AnalysisResult struct {
  Issues []string
  Recommendations []string
}


func AnalyzeCode(lang, code string) (AnalysisResult, error) {
  switch lang {
  case "python":
    return analyzePythonCode(code)
  default:
    return AnalysisResult{}, errors.New("unsupported language")
  }
}

// tmp function with no functionality
func analyzePythonCode(code string) (AnalysisResult, error) {
  fmt.Println("Analyzing python code...")

  issues, recommendations, err := pythonstatic.AnalyzeAST(code)
  if err != nil {
    return AnalysisResult{}, err
  }

  return AnalysisResult{
    Issues: issues,
    Recommendations: recommendations,
  }, nil
}
