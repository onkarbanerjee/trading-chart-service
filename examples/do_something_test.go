package examples_test

import (
	"context"
	"errors"
	"testing"

	"github.com/honestbank/honest-lib-go/hbconfig"
	"github.com/honestbank/honest-lib-go/log"
	"github.com/honestbank/honest-lib-go/tracing"
	"github.com/honestbank/template-repository-go/config"
	"github.com/honestbank/template-repository-go/examples"
)

func TestDoSomethingExample(t *testing.T) {
	traceCfg := hbconfig.From[config.TracesConfig]().GetLatestAvailable()
	closer := tracing.SetTracer("example", traceCfg.Endpoint, traceCfg.Path, log.New())
	defer func() {
		if err := closer(); err != nil {
			t.Log(err)
		}
	}()
	type args struct {
		ctx       func() context.Context
		something func() error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "it should success doing something",
			args: args{
				ctx: func() context.Context {
					_, span := tracing.NewTracer(context.Background(), "example", "test")

					return span.Context()
				},
				something: func() error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "it should error doing something",
			args: args{
				ctx: func() context.Context {
					_, span := tracing.NewTracer(context.Background(), "example", "test")

					return span.Context()
				},
				something: func() error {
					return errors.New("something failed")
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := examples.DoSomethingExample(tt.args.ctx(), tt.args.something); (err != nil) != tt.wantErr {
				t.Errorf("DoSomethingExample() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
