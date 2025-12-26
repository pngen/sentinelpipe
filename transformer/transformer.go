// sentinelpipe/transformer/transformer.go
package transformer

import (
	"context"
)

type Transformer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context)
}

type LogTransformer struct{}

func NewLogTransformer() *LogTransformer {
	return &LogTransformer{}
}

func (t *LogTransformer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (t *LogTransformer) Stop(ctx context.Context) {
	// Cleanup logic
}