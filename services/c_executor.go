package services

import (
	"fmt"
	"time"

	"github.com/eugene817/Cowdocs/container"
)

func makeConfigC(CCode string) container.ContainerConfig {
    script := fmt.Sprintf(`set -eu
cat << 'EOF' > /tmp/main.c
%s
EOF
gcc /tmp/main.c -o /tmp/main.out
/tmp/main.out
`, CCode)
  return container.ContainerConfig{
    Image: "gcc:4.9",
    Cmd: []string{"sh", "-c",  script},
    Tty: false,
  }
}

func (s *Service) ExecuteCInContainer(CCode string) (string, error) {
	contanerConfig := makeConfigC(CCode)
	result, _, err := s.apiSvc.RunContainer(contanerConfig, false)
	if err != nil {
		return "", fmt.Errorf("failed to run container: %v", err)
	}

	return result, nil
}

func (s *Service) ExecuteCWithMetrics(CCode string) (string, string, error) {
	defer timeTrack(time.Now(), "ExecuteCWithMetrics")

	contanerConfig := makeConfigC(CCode) 
	result, metrics, err := s.apiSvc.RunContainer(contanerConfig, false)
	if err != nil {
		return "", "", fmt.Errorf("failed to run container: %v", err)
	}

	return result, metrics, nil
}
