package integration_test

import (
	"context"
	"io"
	"testing"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/onkarbanerjee/trading-chart-service/cmd/aggregator"
	"github.com/stretchr/testify/assert"

	"github.com/onkarbanerjee/trading-chart-service/generated/proto"
	"google.golang.org/grpc"
)

func TestAggregator(t *testing.T) {
	go aggregator.Start()

	// Wait for the server to start
	time.Sleep(2 * time.Second)

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	assert.NoError(t, err)
	defer conn.Close()

	client := proto.NewCandlesClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stream, err := client.Broadcast(ctx, &proto.Empty{})
	assert.NoError(t, err)

	// Read a few candle updates
	received := 0
	for {
		candle, err := stream.Recv()
		if err == io.EOF || ctx.Err() != nil {
			break
		}
		assert.NoError(t, err)
		t.Logf("Received candle: %+v", candle)
		received++
		if received >= 3 {
			break
		}
	}

	assert.Equal(t, 3, received, "should receive 3 candles")
}
