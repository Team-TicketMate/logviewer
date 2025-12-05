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

type flushWriter struct {
	writer  http.ResponseWriter
	flusher http.Flusher
}

// ResponseWriter 가 http.Flusher를 지원하는지 확인하고 감싸주는 함수
func newFlushWriter(writer http.ResponseWriter) (*flushWriter, error) {
	flusher, ok := writer.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("http.ResponseWriter가 스트리밍을 지원하지 않습니다")
	}

	return &flushWriter{
		writer:  writer,
		flusher: flusher,
	}, nil
}

// io.Writer 인터페이스 구현
func (flushWriterInstance *flushWriter) Write(p []byte) (int, error) {
	written, err := flushWriterInstance.writer.Write(p)
	if err != nil {
		return written, err
	}

	// 매번 쓰고 나서 즉시 클라이언트로 flush
	flushWriterInstance.flusher.Flush()
	return written, nil
}

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

	// follow 파라미터 처리: follow = true 또는 follow = 1이면 실시간 스트리밍
	followParam := request.URL.Query().Get("follow")
	follow := strings.EqualFold(followParam, "true") || followParam == "1"

	// 로그는 text/plain 스트리밍 또는 전체 응답
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if follow {
		// 실시간 스트리밍 모드
		streamWriter, newWriterErr := newFlushWriter(writer)
		if newWriterErr != nil {
			log.Printf("HTTP 스트리밍을 지원하지 않는 환경입니다: %v", newWriterErr)
			http.Error(writer, "스트리밍을 지원하지 않는 환경입니다", http.StatusInternalServerError)
			return
		}

		log.Printf("컨테이너 로그 스트리밍 시작 (follow 모드). 컨테이너 ID: %s", containerID)

		// request.Context() 를 넘겨서 클라이언트가 연결을 끊으면 doker 프로세스도 종료
		streamErr := dockercli.StreamContainerLogs(request.Context(), containerID, tailLines, streamWriter)
		if streamErr != nil {
			// 스트리밍 중 에러는 서버 로그로만 출력. 응답 바디에는 보내지 않음
			log.Printf("컨테이너 로그 스트리밍 실패. 컨테이너 ID: %s: %v", containerID, streamErr)
		}
		log.Printf("컨테이너 로그 스트리밍 종료. 컨테이너 ID: %s", containerID)
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
