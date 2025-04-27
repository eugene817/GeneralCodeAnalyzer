package services

import (
	"fmt"
	"time"

	"github.com/eugene817/Cowdocs/container"
)

func makeConfigPython(pythonCode string) container.ContainerConfig {
  return container.ContainerConfig{
    Image: "python:3",
    Cmd: []string{
      "sh", "-c", fmt.Sprintf(`
        echo "%s" > /tmp/script.py &&
        python3 /tmp/script.py
        `, pythonCode),
    },
    Tty: false,
  }
}

func (s *Service) ExecutePythonInContainer(pythonCode string) (string, error) {
	contanerConfig := makeConfigPython(pythonCode)
	result, _, err := s.apiSvc.RunContainer(contanerConfig, false)
	if err != nil {
		return "", fmt.Errorf("failed to run container: %v", err)
	}

	return result, nil
}

func (s *Service) ExecutePythonWithMetrics(pythonCode string) (string, string, error) {
	defer timeTrack(time.Now(), "ExecutePythonWithMetrics")

	contanerConfig := makeConfigPython(pythonCode)
	result, metrics, err := s.apiSvc.RunContainer(contanerConfig, false)
	if err != nil {
		return "", "", fmt.Errorf("failed to run container: %v", err)
	}

	return result, metrics, nil
}
