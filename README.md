# TicketMate Log Viewer

---

## 1. 기본 개념: API 엔드포인트

서비스는 기본적으로 다음 두 가지 API를 제공합니다.

1. **컨테이너 목록 조회**  
   - `GET /containers`
2. **특정 컨테이너 로그 조회**  
   - `GET /containers/{containerId}/logs?lines=500|1000|all&follow=true|false`

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

### 3-2. 최근 1000줄 로그 조회

```bash
curl "https://docker.chuseok22.com/containers/ticketmate-api/logs?lines=1000"
```

### 3-3. 전체 로그 조회

```bash
curl "https://docker.chuseok22.com/containers/ticketmate-api/logs?lines=all"
```

### 3-4. lines 파라미터를 생략하는 경우

`lines` 파라미터를 생략하면 기본값은 `500`으로 동작합니다.

```bash
curl "https://docker.chuseok22.com/containers/ticketmate-api/logs"
```

→ 내부적으로는 `lines=500`과 동일하게 처리됩니다.

---

## 4. 실시간 로그 스트리밍(follow) 기능

`docker logs -f`와 동일한 동작을 HTTP로 제공합니다.

### 4-1. 기본 개념

- 엔드포인트:  
  `GET /containers/{containerId}/logs`
- 지원 쿼리 파라미터:
  - `lines`  
    - `500` / `1000` / `all`  
    - 이전 로그를 얼마나 가져올지 결정
    - 생략 시 기본값은 `500`
  - `follow`  
    - `true` / `1` → **실시간 스트리밍 모드 활성화**
    - 생략 또는 `false` / `0` → 한 번에 응답 반환(기존 동작)

즉,

- `lines` = 과거 로그 몇 줄까지 볼지
- `follow` = 이후에 새로 들어오는 로그를 계속 볼지 여부

라고 이해하면 됩니다.

### 4-2. 최근 500줄 + 이후 로그 실시간 스트리밍

```bash
curl -N "https://docker.chuseok22.com/containers/ticketmate-api/logs?lines=500&follow=true"
```

- `-N` 옵션은 `curl`의 출력 버퍼링을 끄기 위해 사용합니다.  
  (`Flush` 되는 대로 바로바로 로그가 출력되도록 하기 위함입니다.)
- 이미 쌓여 있던 마지막 500줄을 먼저 출력한 뒤,  
  이후에 새로 발생하는 로그를 실시간으로 계속 출력합니다.
- 스트리밍을 중단하려면 `Ctrl + C`로 `curl`을 종료하면 됩니다.

### 4-3. 전체 로그 + 실시간 스트리밍

과거 전체 로그를 모두 출력한 뒤, 이후 로그를 계속 보고 싶다면:

```bash
curl -N "https://docker.chuseok22.com/containers/ticketmate-api/logs?lines=all&follow=true"
```

- 과거 전체 로그를 한 번에 출력한 후,
- 새로 들어오는 로그를 실시간으로 추가 출력합니다.

### 4-4. follow만 사용하고 싶은 경우

`lines`를 생략하고 `follow=true`만 지정하면, 내부적으로는 기본값인 `lines=500`이 적용됩니다.

```bash
curl -N "https://docker.chuseok22.com/containers/ticketmate-api/logs?follow=true"
```

→ `lines=500&follow=true`와 동일하게 동작합니다.

---

## 5. 요약

- `GET /containers`  
  → 현재 동작 중인 컨테이너 목록을 JSON으로 조회
- `GET /containers/{id}/logs?lines=500|1000|all`  
  → 특정 컨테이너의 과거 로그를 한 번에 조회
- `GET /containers/{id}/logs?lines=500|1000|all&follow=true`  
  → 과거 로그 + 이후 들어오는 로그까지 실시간 스트리밍 조회

주로 다음 조합으로 사용하면 됩니다.

1. 컨테이너 목록 확인

   ```bash
   curl "https://docker.chuseok22.com/containers"
   ```

2. 특정 컨테이너 로그 최근 500줄 확인

   ```bash
   curl "https://docker.chuseok22.com/containers/ticketmate-api/logs?lines=500"
   ```

3. 특정 컨테이너 로그 최근 500줄 + 실시간 스트리밍

   ```bash
   curl -N "https://docker.chuseok22.com/containers/ticketmate-api/logs?lines=500&follow=true"
   ```
