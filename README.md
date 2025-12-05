# TicketMate Log Viewer

---

## 1. 기본 개념: API 엔드포인트

서비스는 기본적으로 다음 두 가지 API를 제공합니다.

1. **컨테이너 목록 조회**  
   - `GET /containers`
2. **특정 컨테이너 로그 조회**  
   - `GET /containers/{containerId}/logs?lines=500|1000|all`

여기서 `{containerId}`는 `docker ps` 결과에서 나오는 **컨테이너 ID 또는 이름**입니다.

예시:
- `a1b2c3d4e5f6`
- `ticketmate-api`
- `ticket-mate-logviewer`

---

## 2. curl로 컨테이너 목록 조회하기

### 2-1. 로컬에서 직접 호출

```bash
curl "https://docker.chuseok22.com/containers"
```

응답 예시(JSON):

```json
[
  {
    "id": "a1b2c3d4e5f6",
    "name": "ticketmate-api",
    "image": "ghcr.io/xxx/ticketmate-api:latest",
    "status": "Up 3 hours"
  },
  {
    "id": "7f8e9a0b1c2d",
    "name": "ticket-mate-logviewer",
    "image": "xxx/ticket-mate-logviewer:main",
    "status": "Up 10 minutes"
  }
]
```
---

## 3. curl로 특정 컨테이너 로그 조회하기

### 3-1. 최근 500줄 로그 조회

```bash
curl "https://docker.chuseok22.com/containers/<CONTAINER_ID_OR_NAME>/logs?lines=500"
```

예시:

```bash
curl "https://docker.chuseok22.com/containers/ticketmate-api/logs?lines=500"
```

- `lines=500` → 가장 최근 500줄의 로그를 반환합니다.
- 응답은 `text/plain` 형태의 순수 텍스트입니다.

### 5-2. 최근 1000줄 로그 조회

```bash
curl "https://docker.chuseok22.com/containers/ticketmate-api/logs?lines=1000"
```

### 5-3. 전체 로그 조회

```bash
curl "https://docker.chuseok22.com/containers/ticketmate-api/logs?lines=all"
```

### 5-4. lines 파라미터를 생략하는 경우

`lines` 파라미터를 생략하면 기본값은 `500`으로 동작합니다.

```bash
curl "https://docker.chuseok22.com/containers/ticketmate-api/logs"
```

→ 내부적으로는 `lines=500`과 동일하게 처리됩니다.

