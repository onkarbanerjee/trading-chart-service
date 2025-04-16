package receiver

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/onkarbanerjee/trading-chart-service/entities"
	"go.uber.org/zap"
)

type TicksReceiver struct {
	ctx    context.Context
	wg     *sync.WaitGroup
	conn   *websocket.Conn
	ticks  chan *entities.AggTradeMessage
	logger *zap.Logger
}

func NewTicksReceiver(ctx context.Context, wg *sync.WaitGroup, conn *websocket.Conn, ticks chan *entities.AggTradeMessage, logger *zap.Logger) *TicksReceiver {
	return &TicksReceiver{ctx: ctx, wg: wg, conn: conn, ticks: ticks, logger: logger}
}

func (r *TicksReceiver) Start() {
	defer r.wg.Done()

	operation := "receiver.Start"
	logger := r.logger.With(zap.String("operation", operation))
	logger.Info("listening for trades...")

	defer r.conn.Close()

	for {
		select {
		case <-r.ctx.Done():
			logger.Info("stopping receiver")

			return
		default:
			_, msg, err := r.conn.ReadMessage()
			if err != nil {
				logger.Error("error reading message:", zap.Error(err))

				continue
			}

			var tick entities.AggTradeMessage
			if err := json.Unmarshal(msg, &tick); err != nil {
				logger.Error("error unmarshalling message", zap.Error(err))

				continue
			}

			logger.Info("received tick", zap.Any("tick", tick))
			r.ticks <- &tick
		}
	}
}
