package main

import (
	"log"
	"net/http"
	"ticketmate-logviewer/internal/httpapi"
)

func main() {
	// httpapi 패키지에 라우터 생성을 맡김
	router := httpapi.NewRouter()

	// 서버를 열 포트 지정
	address := ":8080"

	log.Printf("티켓메이트 컨테이너 로그 뷰어 동작 시작: %s", address)

	err := http.ListenAndServe(address, router)
	if err != nil {
		log.Fatalf("서버 실행 실패: %v", err)
	}
}
