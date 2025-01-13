package services

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"
  "runtime"
  "io"
  "encoding/json"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
  "github.com/docker/docker/api/types"
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

func measureMemory() uint64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats.Alloc
}

func ExecuteSQLWithMemoryMetrics(sqlQuery, initSQL string) (string, map[string]interface{}, error) {

	memoryUsedBefore := measureMemory()

	result, metrics, err := ExecuteSQLWithMetrics(sqlQuery, initSQL)
	if err != nil {
		return "", nil, err
	}

	memoryUsedAfter := measureMemory()

	metrics["memory_used_before kiB"] = memoryUsedBefore / 1024
	metrics["memory_used_after kiB"] = memoryUsedAfter / 1024
	metrics["memory_difference kiB"] = (memoryUsedAfter - memoryUsedBefore) / 1024

	return result, metrics, nil
}

func MeasureMemoryUsage() (int, error) {
	query := "PRAGMA cahce_size;"
	result, _, err := ExecuteSQLWithMetrics(query, "")
	if err != nil {
		return 0, fmt.Errorf("failed to measure memory: %w", err)
	}
  log.Printf("Memory used: %s", result)

	var memoryUsed int
	_, err = fmt.Sscanf(result, "%d", &memoryUsed)
	if err != nil {
		return 0, fmt.Errorf("failed to parse memory used: %w", err)
	}

	return memoryUsed, nil
}



func ExecuteSQLMetricsInContainer(sqlQuery string, initSQL string) (string, error) {

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
  
 
  if err := cli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
    return "", fmt.Errorf("failed to start container: %v", err)
  }

  stats, err := cli.ContainerStats(ctx, containerID, true)
	if err != nil {
		panic(err)
	}
	defer stats.Body.Close()
 
		decoder := json.NewDecoder(stats.Body)
		for {
			var stat types.StatsJSON 
			if err := decoder.Decode(&stat); err == io.EOF {
        log.Printf("Error decoding")
				break
			} else if err != nil {
				log.Printf("Error decoding stats: %s\n", err)
				break
			}

			// Печатаем информацию о памяти
			log.Printf("Memory Usage: %v / %v (%.2f%%)\n",
				stat.MemoryStats.Usage,
				stat.MemoryStats.Limit,
				float64(stat.MemoryStats.Usage)/float64(stat.MemoryStats.Limit)*100,
			)

			time.Sleep(1 * time.Second) // Интервал обновления
		}


  defer func() {
    if err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true}); err != nil {
      log.Printf("failed to remove container: %v", err)
    }
  } ()


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
