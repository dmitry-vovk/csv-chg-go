package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dmitry-vovk/csv-chg-go/api"
	"github.com/dmitry-vovk/csv-chg-go/config"
	"github.com/dmitry-vovk/csv-chg-go/source"
	"github.com/dmitry-vovk/csv-chg-go/worker"
)

func main() {
	// Load configuration
	cfg := config.MustLoad()
	// Build worker instance
	w := worker.New(api.New(cfg.APIURL)).
		WithWorkersCount(cfg.Workers).
		WithInterval(cfg.Interval)
	// Read input data
	if err := source.ReadAny(cfg.CSVFile, w.ReadUUIDs); err != nil {
		log.Fatalf("Error reading source file: %s", err)
	}
	// Subscribe to OS signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-signals
		log.Printf("Got %s, initiating shutdown...", s)
		w.Shutdown()
	}()
	// Start the worker
	log.Printf("Worker started")
	w.Run()
	log.Printf("Worker exited")
}
