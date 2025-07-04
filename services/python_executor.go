package services

import (
	"fmt"
	"log"
	"time"

	"github.com/eugene817/Cowdocs/container"
)

func makeConfigPython(pythonCode string) container.ContainerConfig {
	script := fmt.Sprintf(`set -eu
cat << 'EOF' > /tmp/script.py
%s
EOF
exec python3 /tmp/script.py
`, pythonCode)

	return container.ContainerConfig{
		Image: "python:3",
		Cmd:   []string{"sh", "-c", script},
		Tty:   false,
	}
}

func makeConfigPythonLint(code string) container.ContainerConfig {
	script := fmt.Sprintf(`set -eu
cat << 'EOF' > /work/script.py
%s
EOF

flake8 --format=default /work/script.py || true
`, code)

	return container.ContainerConfig{
		Image: "luferchikz/python-flake8:latest",
		Cmd:   []string{"sh", "-c", script},
		Tty:   false,
	}
}

func (s *Service) LintPythonInContainer(code string) (string, error) {
	config := makeConfigPythonLint(code)
	out, _, err := s.apiSvc.RunContainer(config, false)
	if err != nil {
		return "", fmt.Errorf("lint failed: %v", err)
	}

	diagnostics := out
	if diagnostics == "" {
		diagnostics = "No linting issues found."
	}

	return diagnostics, nil
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
	result, metrics, err := s.apiSvc.RunContainer(contanerConfig, true)
	if err != nil {
		return "", "", fmt.Errorf("failed to run container: %v", err)
	}

	log.Println(result, metrics)

	return result, metrics, nil
}
