# event-scope

活動管理與瀏覽平台。

## 技術棧

- Frontend：React + TypeScript
- 前端工具鏈：Vite + Bun
- Backend：Go + Gin
- 後端工具鏈：Air
- Event Streaming：Apache Kafka
- 即時更新：Server-Sent Events（SSE）

## 事件處理架構

本專案預計以 Kafka 處理非同步事件，確保事件在暫時性失敗或重複投遞時，仍能可靠且可追蹤地處理。

```text
React 前端
    ▲
    │ SSE：即時事件結果
    │
Go API / SSE Server
    │
    ├── 發布事件 ──► Kafka Topic ──► Consumer
    │                                  │
    │                                  ├── Idempotency：以事件 ID 去重，避免重複副作用
    │                                  ├── Retry：暫時性錯誤依策略重試
    │                                  └── 成功 ──► SSE 推送結果至前端渲染
    │
    └── 重試耗盡或不可恢復錯誤 ──► DLQ Topic
                                         │
                                         └── 管理端 API／DLQ Consumer 重放至 Retry 或 Main Topic
```

### 處理原則

- **Kafka**：解耦 API 與非同步工作者，並提供可擴充的事件傳遞機制。
- **Idempotency**：所有事件帶有唯一 `event_id`；Consumer 在執行副作用前檢查並記錄處理結果，避免同一事件重複執行。
- **Retry**：可重試的暫時性錯誤採有限次數與退避策略重試。
- **DLQ**：重試耗盡或不可恢復的事件寫入 Kafka 的 Dead Letter Queue Topic，避免阻塞主要事件流。DLQ 事件由管理端 API 或專用 Consumer 重放至 Retry Topic 或原始 Main Topic；不直接自動重呼原始 API，以避免重複驗證或副作用。
- **SSE**：後端將事件處理狀態與結果推送給 React 前端，讓介面可即時渲染進度、成功或失敗狀態。

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
