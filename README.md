# SentinelPipe

A high-performance telemetry pipeline for collecting, transforming, and exporting logs, metrics, and traces.

## Overview

SentinelPipe is a Go-based observability pipeline that ingests telemetry data in multiple formats, processes it through configurable transformers, and exports to various backends including Prometheus, OpenTelemetry, Loki, ClickHouse, and Kafka. Built with performance, scalability, and extensibility in mind.

## Features

- **Multi-format Ingestion**: Supports logs, metrics, and traces
- **Goroutine-Safe Batching**: Efficient concurrent processing with backpressure
- **Plugin Architecture**: Extensible system for collectors, transformers, and sinks
- **Multiple Export Destinations**: Prometheus, OpenTelemetry, Loki, ClickHouse, Kafka
- **Control Plane API**: Chi-based REST API for monitoring and control
- **Full Test Coverage**: Comprehensive unit tests and benchmarks
- **SDK Support**: Embeddable collector SDK for custom integrations

## Architecture

```
[Collectors] → [Transformers] → [Sinks]
     │              │            │
     ▼              ▼            ▼
  Logs        Metrics/Traces   Prometheus
  Metrics         Traces       OpenTelemetry
  Traces         Metrics       Loki
                               ClickHouse
                               Kafka
```

## Getting Started

### Installation

```bash
go get github.com/sentinelpipe/sentinelpipe
```

### Basic Usage

```go
package main

import (
    "context"
    "time"
    
    "github.com/sentinelpipe/sentinelpipe"
    "github.com/sentinelpipe/sentinelpipe/collector"
    "github.com/sentinelpipe/sentinelpipe/transformer"
    "github.com/sentinelpipe/sentinelpipe/sink"
)

func main() {
    // Create pipeline
    pipeline := sentinelpipe.NewPipeline()
    
    // Add components
    pipeline.AddCollector(collector.NewLogCollector(1 * time.Second))
    pipeline.AddTransformer(transformer.NewLogTransformer())
    pipeline.AddSink(sink.NewPrometheusSink())
    
    // Start pipeline
    ctx := context.Background()
    go func() {
        pipeline.StartCollectors(ctx)
        pipeline.StartTransformers(ctx)
        pipeline.StartSinks(ctx)
    }()
    
    // Graceful shutdown
    <-ctx.Done()
    pipeline.Shutdown(ctx)
}
```

### Using the SDK

```go
import "github.com/sentinelpipe/sentinelpipe/sdk"

// Create embedded collector
sdk := sdk.New()
sdk.AddCollector(myCustomCollector)
```

## Components

### Collectors
- **LogCollector**: Collects log data at specified intervals
- Extendable via `collector.Collector` interface

### Transformers
- **LogTransformer**: Processes and transforms telemetry data
- Extendable via `transformer.Transformer` interface

### Sinks
- **PrometheusSink**: Exports metrics to Prometheus
- Extendable via `sink.Sink` interface

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check endpoint |
| `/metrics` | GET | Prometheus metrics endpoint |

## Configuration

The pipeline supports configuration through environment variables and flags. See `config/` directory for examples.

## Performance

- **Concurrency**: Fully concurrent processing with goroutines
- **Backpressure**: Context-based cancellation for graceful shutdowns
- **Memory**: Efficient batching and streaming
- **Scalability**: Horizontal scaling through component composition

## Testing

Run tests with:

```bash
go test -v ./...
```

Run benchmarks with:

```bash
go test -bench=.
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a pull request

## License

MIT License

## Author

Paul Ngen
