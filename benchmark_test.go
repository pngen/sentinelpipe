// sentinelpipe/benchmark_test.go
package main

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/sentinelpipe/sentinelpipe/collector"
)

type benchmarkCollector struct {
	count int64
}

func (c *benchmarkCollector) Start(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			c.count++
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *benchmarkCollector) Stop(ctx context.Context) {}

func BenchmarkPipeline(b *testing.B) {
	pipeline := NewPipeline()
	
	for i := 0; i < b.N; i++ {
		pipeline.AddCollector(&benchmarkCollector{})
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	var wg sync.WaitGroup
	wg.Add(b.N)
	
	for i := 0; i < b.N; i++ {
		go func() {
			defer wg.Done()
			pipeline.StartCollectors(ctx)
		}()
	}
	
	wg.Wait()
}