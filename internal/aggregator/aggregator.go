package aggregator

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/onkarbanerjee/trading-chart-service/entities"
)

type Aggregator interface {
	Aggregate()
}

type OHLCAggregator struct {
	ctx           context.Context
	wg            *sync.WaitGroup
	currentCandle *entities.Candle
	interval      time.Duration
	ticks         <-chan *entities.AggTradeMessage
	candles       chan<- *entities.Candle
	logger        *zap.Logger
}

func NewOHLCAggregator(ctx context.Context, wg *sync.WaitGroup, interval time.Duration, ticks <-chan *entities.AggTradeMessage, candlesChan chan<- *entities.Candle, logger *zap.Logger) *OHLCAggregator {
	return &OHLCAggregator{
		ctx:      ctx,
		wg:       wg,
		interval: interval,
		ticks:    ticks,
		candles:  candlesChan,
		logger:   logger,
	}
}

func (o *OHLCAggregator) Aggregate() {
	defer o.wg.Done()

	for {
		select {
		case <-o.ctx.Done():
			o.logger.Info("aggregator received shutdown signal")

			// Emit the last candle before exiting (if any)
			if o.currentCandle != nil {
				o.candles <- o.currentCandle
			}

			return

		case tick := <-o.ticks:
			bucketTime := time.UnixMilli(tick.Time).Truncate(o.interval)
			log.Println("received tick", zap.Any("tick", tick), zap.Any("bucketTime", bucketTime))
			price, errParsingPrice := strconv.ParseFloat(tick.Price, 64)
			if errParsingPrice != nil {
				o.logger.Error("error parsing price", zap.Error(errParsingPrice))

				continue
			}

			if o.currentCandle == nil || !o.currentCandle.Timestamp.Equal(bucketTime) {
				o.currentCandle = &entities.Candle{
					Symbol:    tick.Symbol,
					Timestamp: bucketTime,
					Open:      price,
					High:      price,
					Low:       price,
					Close:     price,
				}
				o.candles <- o.currentCandle

				continue
			}

			// Update current candle
			o.currentCandle.Close = price
			if price > o.currentCandle.High {
				o.currentCandle.High = price
			}
			if price < o.currentCandle.Low {
				o.currentCandle.Low = price
			}
			o.candles <- o.currentCandle
		}
	}
}
