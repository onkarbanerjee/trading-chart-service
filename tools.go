//go:build tools

//go:generate go run go.uber.org/mock/mockgen -destination=./mocks/mock_grpc_server.go -package=mocks -typed github.com/onkarbanerjee/trading-chart-service/generated/proto Candles_BroadcastServer

package main

import (
	_ "go.uber.org/mock/gomock"
	_ "go.uber.org/mock/mockgen"
	_ "go.uber.org/mock/mockgen/model"
)
