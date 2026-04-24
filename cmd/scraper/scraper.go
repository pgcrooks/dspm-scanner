package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"

	dspm_config "github.com/pgcrooks/dspm-scanner/internal/config"
	"github.com/pgcrooks/dspm-scanner/internal/datastore"
	"github.com/pgcrooks/dspm-scanner/internal/finder"
	"github.com/pgcrooks/dspm-scanner/internal/scanner"
)

func main() {
	handlerOpts := slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				s := a.Value.Any().(*slog.Source)
				s.File = path.Base(s.File)
			}
			return a
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &handlerOpts))
	slog.SetDefault(logger)

	slog.Info("starting scraper orchestrator")

	config, err := dspm_config.GetConfig("config", "./")
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
	finderChan := make(chan finder.BucketObjectBatch, 100)    // finder -> datastore
	dataStoreChan := make(chan finder.BucketObjectBatch, 100) // datastore -> scanner

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

	finderService, err := finder.InitFinderService(ctx, &config, finderChan)
	if err != nil {
		panic(fmt.Errorf("cannot init finder service: %w", err))
	}

	scannerService, err := scanner.InitScannerService(ctx, &config, dataStoreChan)
	if err != nil {
		panic(fmt.Errorf("cannot init scanner service: %w", err))
	}

	// Run workers
	wg.Go(func() {
		datastore.RunDataService(ctx, ds, finderChan)
	})
	wg.Go(func() {
		finderService.Run(ctx)
	})
	wg.Go(func() {
		scannerService.Run(ctx)
	})

	// Run until everything cleans up
	wg.Wait()
	slog.Info("stopped scraper orchestrator")
}
