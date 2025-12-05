package httpapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"ticketmate-logviewer/internal/dockercli"
)

func containersHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "허용되지 않은 HTTP 메서드 압니다", http.StatusMethodNotAllowed)
		return
	}

	containers, err := dockercli.GetRunningContainers()
	if err != nil {
		log.Printf("동작중인 컨테이너 조회에 실패했습니다: %v", err)
		http.Error(writer, "컨테이너 조회에 실패했습니다", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(containers); err != nil {
		log.Printf("컨테이너 응답 인코딩 실패: %v", err)
	}
}

func containerLogsHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "허용되지 않은 HTTP 메서드 입니다", http.StatusMethodNotAllowed)
		return
	}

	path := request.URL.Path
	if !strings.HasPrefix(path, "/containers/") || !strings.HasSuffix(path, "/logs") {
		http.NotFound(writer, request)
		return
	}

	trimmed := strings.TrimPrefix(path, "/containers/")
	containerID := strings.TrimSuffix(trimmed, "/logs")
	containerID = strings.Trim(containerID, "/")

	if containerID == "" {
		http.Error(writer, "컨테이너ID 가 요청되지 않았습니다", http.StatusBadRequest)
		return
	}

	linesParam := request.URL.Query().Get("lines")

	tailLines, err := parseLinesParameter(linesParam)
	if err != nil {
		http.Error(writer, "잘못된 'lines' 파라미터 요청입니다. 500, 1000, all 또는 빈 값만 요청 가능합니다", http.StatusBadRequest)
		return
	}

	logText, err := dockercli.FetchContainerLogs(containerID, tailLines)
	if err != nil {
		log.Printf("컨테이너 로그 조회에 실패했습니다. 컨테이너 ID: %s: %v", containerID, err)
		http.Error(writer, "컨테이너 로그 조회에 실패했습니다", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if _, writeErr := writer.Write([]byte(logText)); writeErr != nil {
		log.Printf("컨테이너 로그 응답 실패: %v", writeErr)
	}
}

func parseLinesParameter(linesParam string) (*int, error) {
	if strings.TrimSpace(linesParam) == "" {
		defaultLines := 500
		return &defaultLines, nil
	}

	switch linesParam {
	case "500", "1000":
		lines, err := strconv.Atoi(linesParam)
		if err != nil {
			return nil, err
		}
		return &lines, nil

	case "all":
		return nil, nil

	default:
		return nil, fmt.Errorf("잘못된 파라미터 요청입니다: %s", linesParam)
	}

}
