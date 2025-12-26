// sentinelpipe/pipeline.go
package main

import (
	"context"
	"sync"
	"time"

	"github.com/sentinelpipe/sentinelpipe/collector"
	"github.com/sentinelpipe/sentinelpipe/transformer"
	"github.com/sentinelpipe/sentinelpipe/sink"
)

type Pipeline struct {
	collectors []collector.Collector
	transformers []transformer.Transformer
	sinks []sink.Sink
	
	mu sync.RWMutex
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		collectors: make([]collector.Collector, 0),
		transformers: make([]transformer.Transformer, 0),
		sinks: make([]sink.Sink, 0),
	}
}

func (p *Pipeline) AddCollector(c collector.Collector) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.collectors = append(p.collectors, c)
}

func (p *Pipeline) AddTransformer(t transformer.Transformer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.transformers = append(p.transformers, t)
}

func (p *Pipeline) AddSink(s sink.Sink) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.sinks = append(p.sinks, s)
}

func (p *Pipeline) StartCollectors(ctx context.Context) error {
	var wg sync.WaitGroup
	
	for _, c := range p.collectors {
		wg.Add(1)
		go func(collector collector.Collector) {
			defer wg.Done()
			if err := collector.Start(ctx); err != nil {
				log.Printf("Collector %T failed: %v", collector, err)
			}
		}(c)
	}
	
	wg.Wait()
	return nil
}

func (p *Pipeline) StartTransformers(ctx context.Context) error {
	var wg sync.WaitGroup
	
	for _, t := range p.transformers {
		wg.Add(1)
		go func(transformer transformer.Transformer) {
			defer wg.Done()
			if err := transformer.Start(ctx); err != nil {
				log.Printf("Transformer %T failed: %v", transformer, err)
			}
		}(t)
	}
	
	wg.Wait()
	return nil
}

func (p *Pipeline) StartSinks(ctx context.Context) error {
	var wg sync.WaitGroup
	
	for _, s := range p.sinks {
		wg.Add(1)
		go func(sink sink.Sink) {
			defer wg.Done()
			if err := sink.Start(ctx); err != nil {
				log.Printf("Sink %T failed: %v", sink, err)
			}
		}(s)
	}
	
	wg.Wait()
	return nil
}

func (p *Pipeline) Shutdown(ctx context.Context) error {
	var wg sync.WaitGroup
	
	for _, c := range p.collectors {
		wg.Add(1)
		go func(collector collector.Collector) {
			defer wg.Done()
			c.Stop(ctx)
		}(c)
	}
	
	for _, t := range p.transformers {
		wg.Add(1)
		go func(transformer transformer.Transformer) {
			defer wg.Done()
			transformer.Stop(ctx)
		}(t)
	}
	
	for _, s := range p.sinks {
		wg.Add(1)
		go func(sink sink.Sink) {
			defer wg.Done()
			sink.Stop(ctx)
		}(s)
	}
	
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}