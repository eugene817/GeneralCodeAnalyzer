package services

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/plugin/executor/containerd"
)


func ExecuteSQLInContainer(sqlQuery string, initSQL string) (string, error) {

  ctx := context.Background()

  cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
  if err != nil {
    return "", fmt.Errorf("failed to create Docker client: %v", err)
  }


  contanerConfig := &container.Config{
    Image: "nouchka/sqlite3",
    Cmd: []string{
      "sh", "-c", fmt.Sprintf(`
        echo "%s" > /tmp/init.sql &&
        sqlite3 /tmp/test.db < /tmp/init.sql &&
        sqlite3 /tmp/test.db "%s"
        `, initSQL, sqlQuery),
    },
    Tty: false,
  }


  resp, err := cli.ConfigCreate(ctx, contanerConfig, nil, nil, nil, "")
  if err != nil {
    return "", fmt.Errorf("failed to create container: %v", err)
  }


  containerID := resp.ID
  defer func() {
    if err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true}); err != nil {
      log.Printf("failed to remove container: %v", err)
    }
  } ()

  
  if err := cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
    return "", fmt.Errorf("failed to start container: %v", err)
  }

  
  logs, err := cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
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

  return stdout.String(), nil
}
