package broadcast

import (
	"github.com/onkarbanerjee/trading-chart-service/entities"
	"github.com/onkarbanerjee/trading-chart-service/generated/proto"
)

type Server struct {
	proto.UnimplementedCandlesServer
	candles <-chan *entities.Candle
}

func NewServer(candles <-chan *entities.Candle) *Server {
	return &Server{candles: candles}
}

func (s *Server) Broadcast(empty *proto.Empty, stream proto.Candles_BroadcastServer) error {
	for candle := range s.candles {
		err := stream.Send(&proto.Candle{
			Symbol:    candle.Symbol,
			Timestamp: candle.Timestamp.UnixMilli(),
			Open:      candle.Open,
			High:      candle.High,
			Low:       candle.Low,
			Close:     candle.Close,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
