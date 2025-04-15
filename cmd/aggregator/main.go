package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/honestbank/template-repository-go/internal/aggregator"

	"github.com/honestbank/template-repository-go/entities"

	"go.uber.org/zap"

	"github.com/gorilla/websocket"
)

func connectToBinance(symbols []string) (*websocket.Conn, error) {
	// Build stream path like btcusdt@aggTrade/ethusdt@aggTrade
	streams := strings.Join(symbols, "@aggTrade/") + "@aggTrade"
	wsURL := fmt.Sprintf("wss://stream.binance.com:9443/stream?streams=%s", streams)

	u, err := url.Parse(wsURL)
	if err != nil {
		return nil, fmt.Errorf("invalid WebSocket URL: %w", err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Binance: %w", err)
	}

	return conn, nil
}

func main() {

	symbols := []string{"btcusdt", "ethusdt", "pepeusdt"}
	conn, err := connectToBinance(symbols)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	ohlcAggregator := aggregator.NewOHLCAggregator(20*time.Second, func(c entities.Candle) {
		logger.Info(fmt.Sprintf("candle: %+v", c))
	})
	logger.Info("Listening for trades...")

	for {
		select {
		case <-interrupt:
			logger.Info("interrupt received. Exiting.")
			conn.Close()

			break
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				logger.Error("error reading message:", zap.Error(err))

				continue
			}

			var tick entities.Tick

			if err := json.Unmarshal(msg, &tick); err != nil {
				logger.Error("error unmarshalling message", zap.Error(err))

				continue
			}

			//logger.Info("received message", zap.Any("message", tick.Data))
			err = ohlcAggregator.Aggregate(&tick)
			if err != nil {
				logger.Error("aggregator error", zap.Error(err))
			}
		}
	}
}
