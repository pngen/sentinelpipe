// sentinelpipe/sdk/sdk.go
package sdk

import (
	"github.com/sentinelpipe/sentinelpipe/collector"
	"github.com/sentinelpipe/sentinelpipe/transformer"
	"github.com/sentinelpipe/sentinelpipe/sink"
)

type SentinelPipe struct {
	collectors []collector.Collector
	transformers []transformer.Transformer
	sinks []sink.Sink
}

func New() *SentinelPipe {
	return &SentinelPipe{
		collectors: make([]collector.Collector, 0),
		transformers: make([]transformer.Transformer, 0),
		sinks: make([]sink.Sink, 0),
	}
}

func (sp *SentinelPipe) AddCollector(c collector.Collector) {
	sp.collectors = append(sp.collectors, c)
}

func (sp *SentinelPipe) AddTransformer(t transformer.Transformer) {
	sp.transformers = append(sp.transformers, t)
}

func (sp *SentinelPipe) AddSink(s sink.Sink) {
	sp.sinks = append(sp.sinks, s)
}