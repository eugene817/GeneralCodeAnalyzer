package services

import (
	"fmt"
	"time"

	"github.com/eugene817/Cowdocs/container"
)

func makeConfigSQL(sqlQuery, initSQL string) container.ContainerConfig {
  return container.ContainerConfig{
    Image: "keinos/sqlite3",
    Cmd: []string{
      "sh", "-c", fmt.Sprintf(`
        echo "%s" > /tmp/init.sql &&
        sqlite3 /tmp/test.db < /tmp/init.sql &&
        sqlite3 /tmp/test.db "%s" > /dev/stdout
        `, initSQL, sqlQuery),
    },
    Tty: false,
  }
}

func (s *Service) ExecuteSQLInContainer(sqlQuery string, initSQL string) (string, error) {

	contanerConfig := makeConfigSQL(sqlQuery, initSQL)
	result, _, err := s.apiSvc.RunContainer(contanerConfig, false)
	if err != nil {
		return "", fmt.Errorf("failed to run container: %v", err)
	}

	return result, nil
}

func (s *Service) ExecuteSQLWithMetrics(sqlQuery, initSQL string) (string, string, error) {
	defer timeTrack(time.Now(), "ExecuteSQLWithMetrics")


	contanerConfig := makeConfigSQL(sqlQuery, initSQL)
	result, metrics, err := s.apiSvc.RunContainer(contanerConfig, true)
	if err != nil {
		return "", "", fmt.Errorf("failed to run container: %v", err)
	}

	return result, metrics, nil
}

func (s *Service) AnalyzeQueryInContainer(sqlQuery, initSQL string) (string, error) {
	explainQuery := fmt.Sprintf("EXPLAIN QUERY PLAN %s;", sqlQuery)
	return s.ExecuteSQLInContainer(explainQuery, initSQL)
}
