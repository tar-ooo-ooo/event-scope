# event-scope

活動管理與瀏覽平台。

## 技術棧

- Frontend：React + TypeScript
- 前端工具鏈：Vite + Bun
- Backend：Go + Gin
- 後端工具鏈：Air
- Event Streaming：Apache Kafka
- 即時更新：Server-Sent Events（SSE）

## 付款事件處理架構

`POST /event` 只負責將付款請求可靠寫入 Kafka，不直接執行付款或推播。背景 worker 主動從 Kafka 拉取事件，讓大量請求可先保存、再依 worker 處理能力執行。

```text
React 前端
    ▲
    │ SSE：accepted／付款／推播結果
    │
Go API / SSE Server
    │
    └── payment.requested ──► Kafka
                                  │
                                  ▼
                    payment-worker consumer group
                    ├── 模擬付款（95% 成功）
                    ├── 暫時性失敗：有限次 retry
                    ├── 重試耗盡：payment.dlq
                    └── 成功：payment.succeeded
                                      │
                                      ▼
                 notification-worker consumer group
                 ├── 模擬推播（95% 成功）
                 ├── 失敗：notification.dlq
                 └── 成功：SSE 推送前端
```

### 處理原則

- **Producer**：API 成功寫入 `payment.requested` 後才回覆 `202 Accepted` 與 SSE `accepted`。
- **Consumer group**：worker 主動拉取 Kafka 訊息；同一 group 中每筆訊息只會由一個 worker 處理。增加 worker 與 partition 後可平行處理。
- **付款優先於推播**：notification worker 只讀取 `payment.succeeded`，付款失敗不會推播。
- **Retry 與 DLQ**：付款重試耗盡寫入 `payment.dlq`；推播失敗寫入 `notification.dlq`。兩者可獨立監控與重放。
- **Idempotency**：所有事件帶有唯一 `event_id`；正式接入真實金流前，consumer 必須以它避免重複副作用。

### DLQ 事件內容

每一筆 DLQ 訊息應保留以下資訊，讓失敗事件可被監控、診斷與安全重放：

- 原始 payload 與唯一 `event_id`
- 原始 Topic、partition 與 offset
- 失敗原因、錯誤碼與錯誤訊息
- 已重試次數與最後失敗時間

## 專案結構

```text
event-scope/
├── frontend/                 # React + TypeScript 前端應用程式
│   ├── public/               # 靜態公開資源
│   ├── src/
│   │   ├── components/       # 可重用 UI 元件
│   │   ├── pages/            # 頁面元件
│   │   ├── services/         # API 與外部服務存取
│   │   ├── types/            # TypeScript 型別定義
│   │   ├── App.tsx           # 應用程式根元件
│   │   └── main.tsx          # 前端進入點
│   ├── package.json
│   ├── bun.lock
│   └── tsconfig.json
├── backend/                  # Go 後端服務
│   ├── .air.toml              # Air 熱重載設定
│   ├── cmd/
│   │   └── api/              # API 服務進入點
│   │       └── main.go
│   ├── internal/
│   │   ├── handler/          # HTTP handlers
│   │   ├── service/          # 商業邏輯
│   │   ├── repository/       # 資料存取層
│   │   ├── model/            # 領域資料模型
│   │   ├── kafka/            # Producer、Consumer 與 Topic 設定
│   │   ├── retry/            # 重試與退避策略
│   │   ├── idempotency/      # 事件去重與處理紀錄
│   │   └── sse/              # SSE 連線與事件推送
│   ├── go.mod
│   └── go.sum
├── README.md
└── .gitignore
```

## 前端開發

前端使用 Bun 管理依賴與執行指令。請先安裝 Bun，再於 `frontend/` 目錄執行：

```bash
bun install
bun run dev
```

常用指令：

- `bun run dev`：啟動本機開發伺服器。
- `bun run build`：執行型別檢查並建立正式版檔案。
- `bun run preview`：預覽正式建置結果。

## 後端開發

後端使用 Gin 建立 HTTP API，並以 Air 在開發期間自動重新編譯與啟動。請先安裝 Go 與 Air，再於 `backend/` 目錄執行：

```bash
air
```

API 預設監聽 `http://localhost:8080`；可透過 `PORT` 環境變數調整連接埠。健康檢查端點為 `GET /healthz`。

## 本機 Kafka

本機使用 [Redpanda](https://redpanda.com/)（與 Kafka 協定相容），並以不含 GUI 的 [Colima](https://colima.run/) 執行 Docker。

首次安裝：

```bash
brew install colima docker
```

每次重新開機後，先啟動 Docker daemon：

```bash
colima start
```

然後在另一個 terminal 前景啟動 broker：

```bash
docker run --rm --name event-scope-kafka -p 9092:9092 \
  docker.redpanda.com/redpandadata/redpanda:latest \
  redpanda start --smp 1 --overprovisioned --node-id 0 \
  --kafka-addr PLAINTEXT://0.0.0.0:9092 \
  --advertise-kafka-addr PLAINTEXT://localhost:9092
```

停止時在該 terminal 按 `Ctrl+C`；容器會自動移除。

建立與查看 Topic（在另一個 terminal 執行）：

```bash
docker exec -it event-scope-kafka rpk topic create payment.requested
docker exec -it event-scope-kafka rpk topic create payment.succeeded
docker exec -it event-scope-kafka rpk topic create payment.dlq
docker exec -it event-scope-kafka rpk topic create notification.dlq
docker exec -it event-scope-kafka rpk topic list
```

Broker 位址為 `localhost:9092`。
