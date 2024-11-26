package services

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)


func ExecuteSQLInContainer(sqlQuery string, initSQL string) (string, error) {

  ctx := context.Background()

  cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
  if err != nil {
    return "", fmt.Errorf("failed to create Docker client: %v", err)
  }


  contanerConfig := &container.Config{
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


  resp, err := cli.ContainerCreate(ctx, contanerConfig, nil, nil, nil, "")
  if err != nil {
    return "", fmt.Errorf("failed to create container: %v", err)
  }

  containerID := resp.ID
  defer func() {
    if err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true}); err != nil {
      log.Printf("failed to remove container: %v", err)
    }
  } ()
  
  if err := cli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
    return "", fmt.Errorf("failed to start container: %v", err)
  }
  
  statusCh, errCh := cli.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
  select {
  case <- statusCh:
  case err := <-errCh:
    return "", fmt.Errorf("error while waiting for container: %w", err)
  }

  logs, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
  if err != nil {
    return "", fmt.Errorf("failed to get logs: %v", err)
  }
  defer logs.Close()

  var stdout, stderr bytes.Buffer
  if _, err := stdcopy.StdCopy(&stdout, &stderr, logs); err != nil {
    return "", fmt.Errorf("failed to read container logs: %v", err)
  }

   if stderr.Len() > 0 {
    return "", fmt.Errorf("error during execution: %s", stderr.String())
  }

  cleanedResult := strings.TrimSpace(stdout.String())
  return cleanedResult, nil
}


func ExecuteSQLWithMetrics(sqlQuery, initSQL string) (string, map[string]interface{}, error) {
  start := time.Now()

  // standart execute
  result, err := ExecuteSQLInContainer(sqlQuery, initSQL) 
  if err != nil {
    return "", nil , err
  }

  elapsed := time.Since(start).Seconds()

  metrics := map[string]interface{} {
    "execution_time": elapsed,
  }

  return result, metrics, nil
}


func AnalyzeQueryInContainer(sqlQuery, initSQL string) (string, error) {
  explainQuery := fmt.Sprintf("EXPLAIN QUERY PLAN %s;", sqlQuery)
  return ExecuteSQLInContainer(explainQuery, initSQL)
}


func ExecuteSQLWithMemoryMetrics(sqlQuery, initSQL string) (string, map[string]interface{}, error) {
	// Выполнение initSQL и запроса
	result, metrics, err := ExecuteSQLWithMetrics(sqlQuery, initSQL)
	if err != nil {
		return "", nil, err
	}

	// Измеряем использование памяти
	memoryUsedBefore, err := MeasureMemoryUsage()
	if err != nil {
		log.Printf("Error measuring memory before: %v", err)
	}

	// Выполняем запрос ещё раз для подсчёта памяти после
	memoryUsedAfter, err := MeasureMemoryUsage()
	if err != nil {
		log.Printf("Error measuring memory after: %v", err)
	}

	// Добавляем метрики памяти
	metrics["memory_used_before"] = memoryUsedBefore
	metrics["memory_used_after"] = memoryUsedAfter
	metrics["memory_difference"] = memoryUsedAfter - memoryUsedBefore

	return result, metrics, nil
}

// MeasureMemoryUsage измеряет текущее использование памяти SQLite
func MeasureMemoryUsage() (int, error) {
	query := "PRAGMA memory_used;"
	result, _, err := ExecuteSQLWithMetrics(query, "")
	if err != nil {
		return 0, fmt.Errorf("failed to measure memory: %w", err)
	}

	var memoryUsed int
	_, err = fmt.Sscanf(result, "%d", &memoryUsed)
	if err != nil {
		return 0, fmt.Errorf("failed to parse memory used: %w", err)
	}

	return memoryUsed, nil
}
