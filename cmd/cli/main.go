package main

import (
	"context"

	"go.uber.org/zap"

	"github.com/honestbank/honest-lib-go/hbconfig"
	"github.com/honestbank/honest-lib-go/log"
	"github.com/honestbank/honest-lib-go/tracing"
	"github.com/honestbank/template-repository-go/config"
	"github.com/honestbank/template-repository-go/examples"
)

func main() {
	traceCfg := hbconfig.From[config.TracesConfig]().GetLatestAvailable()
	closer := tracing.SetTracer("example", traceCfg.Endpoint, traceCfg.Path, log.New())
	defer func() {
		if err := closer(); err != nil {
			log.Info("closing tracer", zap.Error(err))
		}
	}()
	_, span := tracing.NewTracer(context.Background(), "example", "main")

	err := examples.DoSomethingExample(span.Context(), func() error {
		return nil
	})

	if err != nil {
		log.Error("do something example failed", zap.Error(err))
	}
}
