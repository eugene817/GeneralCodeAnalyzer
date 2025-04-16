package services

import (
	"fmt"
	"time"

	"github.com/eugene817/Cowdocs/api"
	"github.com/eugene817/Cowdocs/container"
)

func ExecutePythonInContainer(pythonCode string) (string, error) {
	mgr, err := container.NewDockerManager()
	if err != nil {
		return "", fmt.Errorf("failed to create Docker manager: %v", err)
	}

	mng := api.NewAPI(mgr)

	contanerConfig := container.ContainerConfig{
		Image: "python:3",
		Cmd: []string{
			"sh", "-c", fmt.Sprintf(`
        echo "%s" > /tmp/script.py &&
        python3 /tmp/script.py
        `, pythonCode),
		},
		Tty: false,
	}

	result, _, err := mng.RunContainer(contanerConfig, false)
	if err != nil {
		return "", fmt.Errorf("failed to run container: %v", err)
	}

	return result, nil
}

func ExecutePythonWithMetrics(pythonCode string) (string, string, error) {
	defer timeTrack(time.Now(), "ExecutePythonWithMetrics")

	mgr, err := container.NewDockerManager()
	if err != nil {
		return "", "", fmt.Errorf("failed to create Docker manager: %v", err)
	}

	mng := api.NewAPI(mgr)

	contanerConfig := container.ContainerConfig{
		Image: "python:3",
		Cmd: []string{
			"sh", "-c", fmt.Sprintf(`
        echo "%s" > /tmp/script.py &&
        python3 /tmp/script.py
        `, pythonCode),
		},
		Tty: false,
	}

	result, metrics, err := mng.RunContainer(contanerConfig, false)
	if err != nil {
		return "", "", fmt.Errorf("failed to run container: %v", err)
	}

	return result, metrics, nil
}
