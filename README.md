# Go API Playground

一個功能完整的 Golang API server，使用 Gin framework 和 GORM，提供 JSON 資料抓取和儲存功能，並整合完整的可觀測性（Observability）堆疊。

## 功能特色

- 🚀 使用 Gin framework 建立高效能 API server
- 🗄️ 使用 GORM 與 PostgreSQL 資料庫整合
- 📥 從指定 URL 抓取 JSON 資料
- ✅ JSON 格式驗證
- 💾 將 JSON 資料儲存到資料庫
- 🐳 完整的 Docker 與 docker-compose 支援
- 🔧 豐富的 Makefile 指令
- 📊 **完整的可觀測性堆疊 (OpenTelemetry)**
  - 📈 Prometheus - 指標收集與監控
  - 📝 Loki - 日誌聚合與查詢
  - 🔍 Tempo - 分散式追蹤
  - 📊 Grafana - 視覺化儀表板
  - 🔄 OpenTelemetry Collector - 遙測資料收集器

## API 端點

### GET /health

健康檢查端點

```bash
curl http://localhost:8080/health
```

### GET /api/v1/fetch-json?link={url}

從指定 URL 抓取 JSON 資料並儲存到資料庫

```bash
curl "http://localhost:8080/api/v1/fetch-json?link=https://jsonplaceholder.typicode.com/posts/1"
```

### GET /api/v1/records

取得已儲存的 JSON 記錄（最新 10 筆）

```bash
curl http://localhost:8080/api/v1/records
```

### GET /

API 資訊端點

```bash
curl http://localhost:8080/
```

## 快速開始

### 方法 1: 使用 Docker Compose（推薦）

1. 啟動開發環境：

```bash
make dev-up
```

1. 測試 API：

```bash
make test-api
```

3. 測試 JSON 抓取功能：

```bash
make test-fetch-json
```

4. 停止環境：

```bash
make dev-down
```

### 方法 2: 本地開發

1. 安裝依賴：

```bash
make deps
```

2. 啟動 PostgreSQL（使用 Docker）：

```bash
docker run -d \
  --name postgres-dev \
  -e POSTGRES_DB=jsonapi \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:15-alpine
```

3. 設定環境變數：

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=jsonapi
```

4. 執行應用程式：

```bash
make run
```

## 專案結構

```text
go-playground/
├── main.go                      # API server 主程式
├── models/                      # 資料模型
│   └── json_data.go            # JSON 資料結構定義
├── handlers/                    # API 處理器
│   └── link_handler.go         # link API 處理邏輯
├── database/                    # 資料庫相關
│   └── postgres.go             # PostgreSQL 連接與操作 (GORM)
├── utils/                       # 工具函數
│   └── json_validator.go       # JSON 驗證工具
├── scripts/                     # 測試腳本
│   ├── test-github-api.sh      # GitHub API 測試腳本
│   └── curl-github-meta.sh     # GitHub Meta API 簡單測試
├── otel/                        # OpenTelemetry 觀測性配置
│   ├── Makefile                # OpenTelemetry 相關指令
│   ├── example/                # 示例配置
│   │   ├── base/               # 基礎配置
│   │   └── base-with-otel-collector/ # 包含 OTEL Collector 配置
│   ├── src/                    # 服務配置檔案
│   │   ├── grafana/            # Grafana 配置與儀表板
│   │   ├── loki/               # Loki 日誌配置
│   │   ├── otel-collector/     # OpenTelemetry Collector 配置
│   │   ├── prometheus/         # Prometheus 監控配置
│   │   └── tempo/              # Tempo 追蹤配置
│   └── data/                   # 持久化資料目錄
├── vendor/                      # Go vendor 依賴 (可選)
├── go.mod                       # Go module 定義
├── go.sum                       # 依賴 checksum
├── Dockerfile                   # API server 容器化
├── docker-compose-all.yml       # 完整環境 (API + 資料庫 + 觀測性)
├── docker-compose-go.yml        # 僅 Go API 服務
├── docker-compose-go-otel.yml   # Go API + OpenTelemetry
├── Makefile                     # 建置與管理指令
└── README.md                    # 專案說明文件
```

## Makefile 指令

### 建置與執行

- `make build` - 建置應用程式
- `make run` - 本地執行
- `make deps` - 安裝依賴
- `make clean` - 清理建置檔案

### Docker 指令

- `make docker-build` - 建置 Docker 映像
- `make docker-up` - 啟動 docker-compose
- `make docker-down` - 停止 docker-compose
- `make docker-logs` - 查看日誌

### 開發指令

- `make dev-up` - 啟動開發環境
- `make dev-down` - 停止開發環境
- `make db-reset` - 重置資料庫

### 測試指令

- `make test` - 執行測試
- `make test-coverage` - 執行測試並產生覆蓋率報告
- `make test-api` - 測試基本 API 端點
- `make test-fetch-json` - 測試 JSON 抓取功能
- `make test-get-records` - 測試記錄查詢功能
- `make test-github-meta` - 使用 GitHub Meta API 進行完整測試
- `make curl-github-meta` - 使用 GitHub Meta API 進行簡單測試

## 環境變數

| 變數名稱 | 預設值 | 說明 |
|:---------|:-------|:-----|
| `DB_HOST` | localhost | 資料庫主機位址 |
| `DB_PORT` | 5432 | 資料庫連接埠 |
| `DB_USER` | postgres | 資料庫使用者名稱 |
| `DB_PASSWORD` | postgres | 資料庫密碼 |
| `DB_NAME` | jsonapi | 資料庫名稱 |
| `PORT` | 8080 | API server 連接埠 |
| `GIN_MODE` | debug | Gin 執行模式 (debug/release) |

## 資料庫結構

```sql
CREATE TABLE json_records (
    id SERIAL PRIMARY KEY,
    url VARCHAR(2048) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 技術棧

### 核心技術
- **框架**: Gin (Golang web framework)
- **ORM**: GORM
- **資料庫**: PostgreSQL
- **HTTP Client**: Resty
- **容器化**: Docker & Docker Compose

### 可觀測性技術
- **指標收集**: Prometheus
- **日誌管理**: Loki
- **分散式追蹤**: Tempo
- **視覺化**: Grafana (含多種專用插件)
- **遙測收集**: OpenTelemetry Collector

## 開發說明

### 錯誤處理

- URL 格式驗證
- HTTP 請求失敗處理
- JSON 格式驗證
- 資料庫操作錯誤處理
- 適當的 HTTP 狀態碼回應

### JSON 驗證

- 檢查 Content-Type 是否為 JSON
- 驗證回應資料是否為有效的 JSON 格式
- 支援 `application/json` 和 `text/json` 類型

### 安全性

- URL 格式驗證（僅允許 HTTP/HTTPS）
- CORS 支援
- SQL Injection 防護（使用 GORM）

## 範例使用

1. 啟動服務：

   ```bash
   make dev-up
   ```

2. 抓取 JSON 資料：

   ```bash
   curl "http://localhost:8080/api/v1/fetch-json?link=https://api.github.com/users/octocat"
   ```

3. 查看儲存的記錄：

   ```bash
   curl http://localhost:8080/api/v1/records
   ```

## 可觀測性功能

本專案整合了完整的 OpenTelemetry 可觀測性堆疊，提供指標監控、日誌管理和分散式追蹤功能。

### 服務端點

當啟動完整環境後，可以透過以下端點存取各個服務：

- **API Server**: http://localhost:8080
- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Loki**: http://localhost:3100

### Grafana 儀表板

專案包含多個預配置的 Grafana 儀表板：

- **RED 方法儀表板**: 監控請求速率、錯誤率和持續時間
- **Go 進程監控**: 監控 Go 應用程式的執行時指標
- **分析資訊儀表板**: 業務指標和使用分析
- **日誌分析儀表板**: 日誌聚合和查詢
- **Pond Pool 資訊**: 自定義業務指標

### OpenTelemetry 配置

專案提供多種 OpenTelemetry 配置選項：

- **基礎配置** (`otel/example/base/`): 最小化的觀測性設定
- **完整配置** (`otel/example/base-with-otel-collector/`): 包含 OTEL Collector 的完整設定

### 啟動觀測性環境

```bash
# 啟動包含觀測性功能的完整環境
make dev-up

# 或使用 OpenTelemetry 特定的 docker-compose 檔案
cd otel/example/base
docker-compose up -d
```

## 故障排除

### 常見問題

1. **資料庫連接失敗**: 確保 PostgreSQL 已啟動且連線資訊正確
2. **依賴下載失敗**: 執行 `make deps` 更新依賴
3. **Docker 權限問題**: 確保用戶已加入 docker 群組
4. **Grafana 插件問題**: 重啟 Grafana 容器或檢查插件安裝狀態
5. **觀測性資料遺失**: 檢查 OpenTelemetry Collector 配置和服務連接狀態

### 日誌查看

```bash
# 查看所有服務日誌
make docker-logs

# 查看特定服務日誌
docker-compose -f docker-compose-all.yml logs -f [服務名稱]
```

### 服務健康檢查

```bash
# 檢查 API 服務狀態
curl http://localhost:8080/health

# 檢查 Prometheus 指標端點
curl http://localhost:9090/metrics

# 檢查 Grafana 狀態
curl http://localhost:3000/api/health
```
