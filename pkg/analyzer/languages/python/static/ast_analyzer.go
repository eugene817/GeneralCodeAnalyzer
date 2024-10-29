package pythonstatic

import (
  "context"
  "encoding/json"
  "io"
  "github.com/docker/docker/api/types"
  "github.com/docker/docker/api/types/container"
  "github.com/docker/docker/client"
	"fmt"
)

func AnalyzeAST(code string) ([]string, []string, error) {
  cli, err := client.NewClientWithOpts(client.FromEnv)
  if err != nil {
    return nil, nil, fmt.Errorf("error creating docker-client: %v", err)
  }

  ctx := context.Background()

  containerConfig := &container.Config{
    Image: "python-analyzer",
    Cmd: []string{"python3", "/app/analyze_script.py"},
    Tty: false,
    OpenStdin: true,
    AttachStdin: true,
    AttachStdout: true,
    AttachStderr: true,
  }

  resp, err := cli.ConfigCreate(ctx, containerConfig, nil, nil, nil, "")
  if err != nil {
    return nil, nil, fmt.Errorf("error creating container: %v", err)
  }

  if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
    return nil, nil, fmt.Errorf("error running container %v", err)
  }

  writer, err := cli.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
    Stream: true,
    Stdin: true,
    Stdout: true,
    Stderr: true,
  })

  if err != nil {
    return nil, nil, fmt.Errorf("error connecting to the container: %v", err)
  }

  defer writer.Close()

  _, err = writer.Conn.Write([]byte(code))
  if err != nil {
    return nil, nil, fmt.Errorf("error writing code to the container: %v", err)
  }

  reader, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
  
  if err != nil {
    return nil, nil, fmt.Errorf("error reading container log: %v", err)
  }

  defer reader.Close()

  output, err := io.ReadAll(reader)
  if err != nil {
    return nil, nil, fmt.Errorf("error reading container output: %v", err)
  }

  var result struct {
    Issues []string `json:"issues"`
    Recommendations []string `json:"recommendations"`
  }

  if err := json.Unmarshal(output, &result); err != nil {
    return nil, nil, fmt.Errorf("error parsing JSON: %v", err)
  }

  cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true})

  return result.Issues, result.Recommendations, nil
}

