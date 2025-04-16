package services

import (
	"strings"
)

func GenerateRecommendations(metrics, analysis, sqlQuery string) []string {

	recommendations := []string{}

	if !contains(analysis, "USING INDEX") {
		recommendations = append(recommendations, "No indexes are used in this query. Consider adding indexes.")
	}

	if strings.Contains(strings.ToUpper(sqlQuery), "SELECT *") {
		recommendations = append(recommendations, "Avoid using SELECT *; specify the columns explicitly.")
	}

	return recommendations
}

func GenerateRecommendationsPython(metrics, pythonCode string) []string {

	recommendations := []string{}

	if strings.Contains(strings.ToUpper(pythonCode), "SELECT *") {
		recommendations = append(recommendations, "Avoid using SELECT *; specify the columns explicitly.")
	}

	return recommendations
}

func contains(s, substr string) bool {
	return strings.Contains(strings.ToUpper(s), strings.ToUpper(substr))
}
