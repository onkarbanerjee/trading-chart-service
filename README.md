# Trading Chart Service

A Go-based trading chart service that ingests real-time tick data from Binance, aggregates it into 1-minute OHLC candlesticks, streams current candle updates via gRPC, and persists completed candles to a datastore. Designed for deployment on Kubernetes using Terraform.

---

## 🔧 Technologies

- Go
- gRPC
- Binance WebSocket API
- PostgreSQL / In-memory store (toggle)
- Kubernetes
- Terraform

---

## 🚀 Getting Started

### Prerequisites

- Go 1.20+
- Docker
- kubectl
- Terraform
- `protoc` compiler with Go plugins

### 1. Clone and Build

```bash
git clone https://github.com/yourusername/trading-chart-service.git
cd trading-chart-service
go build ./cmd/aggregator
```

### 2. Run Locally
```bash
go run ./cmd/aggregator
```

This will start:
- The Binance WebSocket client
- The candle aggregator
- The gRPC streaming server

### 3. Run Tests
```bash
go integration_test ./internal/aggregator
```

### 📦 Kubernetes Deployment

1. Build Docker Image

```bash
docker build -t trading-chart-service:latest .
```
