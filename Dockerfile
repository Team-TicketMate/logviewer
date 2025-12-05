# Go 빌드
FROM golang:1.22-alpine AS builer

# 작업 디렉토리 설정
WORKDIR /app

# Go 모듈 파일 복사 (캐시)
COPY go.mod go.sum ./
RUN go mod download

# 나머지 소스 복사
COPY . .

# CGO 비활성화, Linux/amd64 타겟 빌드
# ./cmd/logviewer/main.go 를 진입점으로 logviwer 바이너리 생성
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o logviewer ./cmd/logviewer

# 런타임 이미지
FROM alpine:3.20

# 앱 실행용 비루트 유저 생성
RUN adduser -D -H -u 10001 appuser \
  && apk add --no-cache ca-certificates docker-cli

WORKDIR /app

# 빌더 단계에서 만든 Go 바이너리 복사
COPY --from=builder /app/logviewer .

# 8080포트
EXPOSE 8080

# 비루트 유저로 실행
USER appuser

# 컨테이너 시작시 실행할 명령
ENTRYPOINT ["./logviewer"]