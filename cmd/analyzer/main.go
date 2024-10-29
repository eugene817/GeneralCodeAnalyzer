package main

import (
  "fmt"
  "github.com/eugene817/GeneralCodeAnalyzer/pkg/analyzer"
)

func main() {
  fmt.Println("GeneralCodeAnalyzer is running...\nWelcome to Code Analyzer")

  code := `
  def factorial(n):
    if n == 1:
      return 1
    else:
      return n * factorial(n-1)
  `

  result, err := analyzer.AnalyzeCode("python", code)
  if err != nil {
    fmt.Println("Error during analysis", err)
    return
  }

  fmt.Println("Analysis Result:", result)
}

