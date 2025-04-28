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
gcc -std=c99 -O2 /tmp/main.c -o /tmp/main.out
exec /tmp/main.out
`, CCode)

    return container.ContainerConfig{
        Image: "gcc:4.9",
        Cmd:   []string{"sh", "-c", script},
        Tty:   false,
    }
}


func makeConfigCLint(code string) container.ContainerConfig {
	script := fmt.Sprintf(`set -eu
cat << 'EOF' > /tmp/main.c
%s
EOF
clang -fsyntax-only /tmp/main.c
`, code)
	return container.ContainerConfig{
		Image: "silkeh/clang", 		
    Cmd:   []string{"sh", "-c", script},
		Tty:   false,
	}
}

func (s *Service) LintCInContainer(code string) (string, error) {
	config := makeConfigCLint(code)
	out, _, err := s.apiSvc.RunContainer(config, false)
	if err != nil {
		return "", fmt.Errorf("lint failed: %v", err)
	}
	return out, nil
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
