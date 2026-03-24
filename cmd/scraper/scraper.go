package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

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

	ds, err := scanner_int.InitDataStore(context.TODO(), scanner_int.LocalDB)
	if err != nil {
		panic(fmt.Errorf("cannot init ds: %w", err))
	}
	defer ds.Close()

	// Communication channels
	scrapeMessages := make(chan []scanner_int.BucketObject, 100)

	var wg sync.WaitGroup
	wg.Go(func() {
		ds.RunDataService(context.TODO(), scrapeMessages)
	})

	if config.Aws.Enabled {
		slog.Debug("aws enabled")

		client, err := newS3Client()
		if err != nil {
			slog.Error("unable to create AWS client", "err", err.Error())
		} else {
			contents, err := scanner_int.ListS3Bucket(
				context.TODO(), client, config.Aws.BucketName,
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

	if config.Local.Enabled {
		slog.Info("local enabled")

		wg.Go(func() {
			scanner_int.RunScraperService(context.TODO(), config, scrapeMessages)
		})
	} else {
		slog.Debug("local disabled")
	}

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
