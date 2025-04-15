package aggregator

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/honestbank/template-repository-go/entities"
)

type Aggregator interface {
	Aggregate(tick *entities.Tick) error
}

type OHLCAggregator struct {
	mu       *sync.RWMutex
	candles  map[string]*entities.Candle
	interval time.Duration
	OnClose  func(candle entities.Candle) // Called when a candle closes
}

func NewOHLCAggregator(interval time.Duration, onClose func(candle entities.Candle)) *OHLCAggregator {
	return &OHLCAggregator{
		mu:       &sync.RWMutex{},
		candles:  make(map[string]*entities.Candle),
		interval: interval,
		OnClose:  onClose,
	}
}

func (o *OHLCAggregator) Aggregate(tick *entities.Tick) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	bucketTime := time.UnixMilli(tick.Data.Time).Truncate(o.interval)
	key := tick.Data.Symbol
	price, errParsingPrice := strconv.ParseFloat(tick.Data.Price, 64)
	if errParsingPrice != nil {
		return fmt.Errorf("could not parse the price, got error: %s", errParsingPrice.Error())
	}

	c, exists := o.candles[key]
	if !exists || !c.Timestamp.Equal(bucketTime) {
		if exists {
			o.OnClose(*c)
		}

		// Create a new candle for this bucket
		o.candles[key] = &entities.Candle{
			Symbol:    tick.Data.Symbol,
			Timestamp: bucketTime,
			Open:      price,
			High:      price,
			Low:       price,
			Close:     price,
		}

		return nil
	}

	// Update current candle
	c.Close = price
	if price > c.High {
		c.High = price
	}
	if price < c.Low {
		c.Low = price
	}

	return nil
}
