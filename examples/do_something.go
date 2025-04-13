package examples

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/honestbank/honest-lib-go/metric"
	"github.com/honestbank/honest-lib-go/tracing"
)

func DoSomethingExample(ctx context.Context, something func() error) error {
	const operation = "DoSomethingExample"
	span := tracing.FromContext(ctx).StartSpan(operation)
	span.Logger().Info("doing something")
	m := metric.NewMetrics()
	measurement := m.StartOperationTimer(operation)
	time.Sleep(2 * time.Second)
	measurement.RecordSuccess()
	if err := something(); err != nil {
		span.Logger().Error("do something failed", zap.Error(err))

		return err
	}

	span.Logger().Info("success doing something")

	return nil
}
