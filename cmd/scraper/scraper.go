package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	scanner_int "github.com/pgcrooks/dspm-scanner/internal"
	datastore "github.com/pgcrooks/dspm-scanner/internal/datastore"
	finder "github.com/pgcrooks/dspm-scanner/internal/finder"
)

func main() {
	slog.Info("starting scraper orchestrator")

	config, err := scanner_int.GetConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Contexts
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// DataStore
	ds, err := datastore.InitDataStore(ctx, &config)
	if err != nil {
		panic(fmt.Errorf("cannot init ds: %w", err))
	}
	defer ds.Close()

	// Communication channels
	scrapeChan := make(chan finder.BucketObjectBatch, 100)

	// Group all routines
	wg := sync.WaitGroup{}

	// Handle shutdown signals
	go func() {
		interruptChan := make(chan os.Signal, 1)
		signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
		<-interruptChan

		slog.Info("received shutdown signal, stopping gracefully")

		cancel()
	}()

	finderService, err := finder.InitFinderService(ctx, &config, scrapeChan)
	if err != nil {
		panic(fmt.Errorf("cannot init finder service: %w", err))
	}

	// Run workers
	wg.Go(func() {
		datastore.RunDataService(ctx, ds, scrapeChan)
	})
	wg.Go(func() {
		finderService.Run(ctx)
	})

	// Run until everything cleans up
	wg.Wait()
	slog.Info("stopped scraper orchestrator")
}
