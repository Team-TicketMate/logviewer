package dockercli

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Container struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

func GetRunningContainers() ([]Container, error) {
	cmd := exec.Command(
		"docker",
		"ps",
		"--format",
		"{{.ID}}||{{.Names}}||{{.Image}}||{{.Status}}",
	)

	outputBytes, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("'docker ps' 실행에 실패했습니다: %w", err)
	}

	output := strings.TrimSpace(string(outputBytes))

	if output == "" {
		return []Container{}, nil
	}

	lines := strings.Split(output, "\n")
	containers := make([]Container, 0, len(lines))

	for _, line := range lines {
		parts := strings.Split(line, "||")
		if len(parts) != 4 {
			continue
		}
		container := Container{
			ID:     strings.TrimSpace(parts[0]),
			Name:   strings.TrimSpace(parts[1]),
			Image:  strings.TrimSpace(parts[2]),
			Status: strings.TrimSpace(parts[3]),
		}
		containers = append(containers, container)
	}

	return containers, nil
}

func FetchContainerLogs(containerID string, tailLines *int) (string, error) {
	trimmedID := strings.TrimSpace(containerID)
	if trimmedID == "" {
		return "", fmt.Errorf("컨테이너 ID가 비어있습니다")
	}

	args := []string{"logs"}

	if tailLines != nil {
		if *tailLines <= 0 {
			return "", fmt.Errorf("tailLines 는 0 이상의 숫자여야합니다")
		}
		args = append(args, "--tail", fmt.Sprintf("%d", *tailLines))
	}

	args = append(args, trimmedID)

	cmd := exec.Command("docker", args...)

	var stdoutStderr bytes.Buffer
	cmd.Stdout = &stdoutStderr
	cmd.Stderr = &stdoutStderr

	err := cmd.Run()
	output := stdoutStderr.String()

	if err != nil {
		return "", fmt.Errorf("'docker logs' 실행에 실패했습니다: %w, 출력: %s", err, output)
	}

	return output, nil
}
