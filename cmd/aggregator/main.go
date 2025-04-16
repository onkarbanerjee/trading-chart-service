package aggregator

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/onkarbanerjee/trading-chart-service/entities"
	"github.com/onkarbanerjee/trading-chart-service/generated/proto"
	"github.com/onkarbanerjee/trading-chart-service/internal/aggregator"
	broadcast "github.com/onkarbanerjee/trading-chart-service/internal/grpc"
	"github.com/onkarbanerjee/trading-chart-service/internal/receiver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func connectToBinance(symbol string) (*websocket.Conn, error) {
	streams := symbol + "@aggTrade"
	wsURL := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s", streams)

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
	Start()
}

func Start() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	candles := make(chan *entities.Candle)
	stop := startGRPCBroadcastingServer(broadcast.NewServer(candles), lis)
	symbols := []string{"btcusdt", "ethusdt", "pepeusdt"}
	for _, symbol := range symbols {
		conn, err := connectToBinance(symbol)
		if err != nil {
			logger.Error("could not get a connection to binance, will check for next symbol", zap.String("symbol", symbol), zap.Error(err))

			continue
		}
		ticks := make(chan *entities.AggTradeMessage)
		wg.Add(1)
		go aggregator.NewOHLCAggregator(ctx, wg, 30*time.Second, ticks, candles, logger.With(zap.String("symbol", symbol))).Aggregate()
		wg.Add(1)
		go receiver.NewTicksReceiver(ctx, wg, conn, ticks, logger.With(zap.String("symbol", symbol))).Start()
	}

	<-interrupt
	logger.Info("interrupt received")
	cancel()
	logger.Info("all contexts cancelled")
	wg.Wait()
	logger.Info("all goroutines finished")
	stop()
	logger.Info("stopped the grpc server")
}

func startGRPCBroadcastingServer(s proto.CandlesServer, lis net.Listener) func() {
	grpcServer := grpc.NewServer()
	proto.RegisterCandlesServer(grpcServer, s)

	log.Println("🚀 gRPC server listening on :50051")
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return grpcServer.GracefulStop
}
