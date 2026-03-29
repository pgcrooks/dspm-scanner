package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	scanner_int "github.com/pgcrooks/dspm-scanner/internal"
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
	ds, err := scanner_int.InitDataStore(ctx, &config)
	if err != nil {
		panic(fmt.Errorf("cannot init ds: %w", err))
	}
	defer ds.Close()

	// Communication channels
	scrapeMessages := make(chan scanner_int.BucketObjectBatch, 100)

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

	// Run workers
	wg.Go(func() {
		ds.Run(ctx, scrapeMessages)
	})

	if config.Scraper.Aws.Enabled {
		slog.Debug("aws enabled")

		client, err := newS3Client()
		if err != nil {
			slog.Error("unable to create AWS client", "err", err.Error())
		} else {
			contents, err := scanner_int.ListS3Bucket(
				context.TODO(), client, config.Scraper.Aws.BucketName,
			)
			if err != nil {
				slog.Error(err.Error())
			} else {
				slog.Info("first page results")
				for _, object := range contents {
					slog.Info("key=%s size=%d", object.Key, object.Size)
				}
			}
		}

	} else {
		slog.Debug("aws disabled")
	}

	if config.Scraper.Local.Enabled {
		slog.Info("local enabled")

		wg.Go(func() {
			scanner_int.RunScraperService(ctx, config, scrapeMessages)
		})
	} else {
		slog.Debug("local disabled")
	}

	// Run until everything cleans up
	wg.Wait()
	slog.Info("stopped scraper orchestrator")
}

func newS3Client() (*s3.Client, error) {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	return client, nil
}
