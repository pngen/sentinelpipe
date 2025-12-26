// sentinelpipe/pipeline_test.go
package main

import (
	"context"
	"testing"
	"time"

	"github.com/sentinelpipe/sentinelpipe/collector"
	"github.com/sentinelpipe/sentinelpipe/transformer"
	"github.com/sentinelpipe/sentinelpipe/sink"
)

type mockCollector struct{}

func (c *mockCollector) Start(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func (c *mockCollector) Stop(ctx context.Context) {}

type mockTransformer struct{}

func (t *mockTransformer) Start(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func (t *mockTransformer) Stop(ctx context.Context) {}

type mockSink struct{}

func (s *mockSink) Start(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func (s *mockSink) Stop(ctx context.Context) {}

func TestPipeline(t *testing.T) {
	pipeline := NewPipeline()
	
	pipeline.AddCollector(&mockCollector{})
	pipeline.AddTransformer(&mockTransformer{})
	pipeline.AddSink(&mockSink{})
	
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	go func() {
		if err := pipeline.StartCollectors(ctx); err != nil {
			t.Logf("Collector error: %v", err)
		}
	}()
	
	go func() {
		if err := pipeline.StartTransformers(ctx); err != nil {
			t.Logf("Transformer error: %v", err)
		}
	}()
	
	go func() {
		if err := pipeline.StartSinks(ctx); err != nil {
			t.Logf("Sink error: %v", err)
		}
	}()
	
	if err := pipeline.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}
}