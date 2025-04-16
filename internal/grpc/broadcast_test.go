package broadcast_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/onkarbanerjee/trading-chart-service/generated/proto"

	"github.com/onkarbanerjee/trading-chart-service/entities"
	broadcast "github.com/onkarbanerjee/trading-chart-service/internal/grpc"
	"github.com/onkarbanerjee/trading-chart-service/mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestServer_Broadcast(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctl := gomock.NewController(t)
		defer ctl.Finish()
		mockStream := mocks.NewMockCandles_BroadcastServer(ctl)

		candles := make(chan *entities.Candle)
		server := broadcast.NewServer(candles)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()

			assert.NoError(t, server.Broadcast(nil, mockStream))
		}()

		now := time.Now()
		expectedTimestamp := now.UnixMilli()
		mockCandles := []*entities.Candle{
			{Symbol: "btcusdt", Timestamp: now, Open: 100, High: 110, Low: 90, Close: 105},
			{Symbol: "ethusdt", Timestamp: now, Open: 200, High: 210, Low: 190, Close: 205},
			{Symbol: "ethusdt", Timestamp: now, Open: 300, High: 310, Low: 290, Close: 305},
			{Symbol: "pepeusdt", Timestamp: now, Open: 400, High: 410, Low: 390, Close: 405},
		}
		lo.ForEach[*entities.Candle](mockCandles, func(item *entities.Candle, index int) {
			mockStream.EXPECT().Send(&proto.Candle{
				Symbol:    item.Symbol,
				Timestamp: expectedTimestamp,
				Open:      item.Open,
				High:      item.High,
				Low:       item.Low,
				Close:     item.Close,
			}).Return(nil)
			candles <- item
		})
		close(candles)
		wg.Wait()
	})
	t.Run("error - from stream", func(t *testing.T) {
		ctl := gomock.NewController(t)
		defer ctl.Finish()
		mockStream := mocks.NewMockCandles_BroadcastServer(ctl)

		candles := make(chan *entities.Candle)
		server := broadcast.NewServer(candles)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()

			assert.ErrorContains(t, server.Broadcast(nil, mockStream), "broadcast error")
		}()

		now := time.Now()
		expectedTimestamp := now.UnixMilli()
		mockCandles := []*entities.Candle{
			{Symbol: "btcusdt", Timestamp: now, Open: 100, High: 110, Low: 90, Close: 105},
			{Symbol: "ethusdt", Timestamp: now, Open: 200, High: 210, Low: 190, Close: 205},
		}
		mockStream.EXPECT().Send(&proto.Candle{
			Symbol:    "btcusdt",
			Timestamp: expectedTimestamp,
			Open:      100,
			High:      110,
			Low:       90,
			Close:     105,
		}).Return(nil)
		mockStream.EXPECT().Send(&proto.Candle{
			Symbol:    "ethusdt",
			Timestamp: expectedTimestamp,
			Open:      200,
			High:      210,
			Low:       190,
			Close:     205,
		}).Return(errors.New("broadcast error"))

		lo.ForEach[*entities.Candle](mockCandles, func(item *entities.Candle, index int) {
			candles <- item
		})
		close(candles)
		wg.Wait()
	})
}
