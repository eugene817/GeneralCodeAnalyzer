package services

import (
  "strings"
)

func GenerateRecommendations(metrics map[string]interface{}, analysis string, sqlQuery string) []string {

	recommendations := []string{}

	if metrics["execution_time"].(float64) > 1.0 {
		recommendations = append(recommendations, "Consider optimizing your query as it takes more than 1 second.")
	}

  if !contains(analysis, "USING INDEX") {
    recommendations = append(recommendations, "No indexes are used in this query. Consider adding indexes.")
  }

	if strings.Contains(strings.ToUpper(sqlQuery), "SELECT *") {
		recommendations = append(recommendations, "Avoid using SELECT *; specify the columns explicitly.")
	}

	return recommendations
}

func contains(s, substr string) bool {
	return strings.Contains(strings.ToUpper(s), strings.ToUpper(substr))
}
