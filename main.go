// sentinelpipe/main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize pipeline
	pipeline := NewPipeline()
	
	// Start collectors
	go func() {
		if err := pipeline.StartCollectors(ctx); err != nil {
			log.Printf("Collector error: %v", err)
		}
	}()

	// Start transformers
	go func() {
		if err := pipeline.StartTransformers(ctx); err != nil {
			log.Printf("Transformer error: %v", err)
		}
	}()

	// Start sinks
	go func() {
		if err := pipeline.StartSinks(ctx); err != nil {
			log.Printf("Sink error: %v", err)
		}
	}()

	// Setup control plane API
	r := chi.NewRouter()
	r.Get("/metrics", promhttp.Handler().ServeHTTP)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	go func() {
		if err := http.ListenAndServe(":8080", r); err != nil {
			log.Printf("API server error: %v", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := pipeline.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}
}