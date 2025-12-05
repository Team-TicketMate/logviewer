package httpapi

import "net/http"

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// 컨테이너 목록 조회: GET /containers
	mux.HandleFunc("/containers", containersHandler)

	// 컨테이너 로그 조회: GET /containers/{id}/logs?lines=500|1000|all
	mux.HandleFunc("/containers/", containerLogsHandler)

	return mux
}
