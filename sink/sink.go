// sentinelpipe/sink/sink.go
package sink

import (
	"context"
)

type Sink interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context)
}

type PrometheusSink struct{}

func NewPrometheusSink() *PrometheusSink {
	return &PrometheusSink{}
}

func (s *PrometheusSink) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *PrometheusSink) Stop(ctx context.Context) {
	// Cleanup logic
}