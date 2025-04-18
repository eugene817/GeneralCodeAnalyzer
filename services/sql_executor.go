package services

import (
	"fmt"
	"time"

	"github.com/eugene817/Cowdocs/api"
	"github.com/eugene817/Cowdocs/container"
)

func ExecuteSQLInContainer(sqlQuery string, initSQL string) (string, error) {

	mgr, err := container.NewDockerManager()
	if err != nil {
		return "", fmt.Errorf("failed to create Docker manager: %v", err)
	}

	mng := api.NewAPI(mgr)

	contanerConfig := container.ContainerConfig{
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

	result, _, err := mng.RunContainer(contanerConfig, false)
	if err != nil {
		return "", fmt.Errorf("failed to run container: %v", err)
	}

	return result, nil
}

func ExecuteSQLWithMetrics(sqlQuery, initSQL string) (string, string, error) {
	defer timeTrack(time.Now(), "ExecuteSQLWithMetrics")

	mgr, err := container.NewDockerManager()
	if err != nil {
		return "", "", fmt.Errorf("failed to create Docker manager: %v", err)
	}

	mng := api.NewAPI(mgr)

	contanerConfig := container.ContainerConfig{
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

	result, metrics, err := mng.RunContainer(contanerConfig, true)
	if err != nil {
		return "", "", fmt.Errorf("failed to run container: %v", err)
	}

	return result, metrics, nil
}

func AnalyzeQueryInContainer(sqlQuery, initSQL string) (string, error) {
	explainQuery := fmt.Sprintf("EXPLAIN QUERY PLAN %s;", sqlQuery)
	return ExecuteSQLInContainer(explainQuery, initSQL)
}
