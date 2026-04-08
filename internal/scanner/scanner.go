package scanner

import (
	"context"
	"log/slog"
	"sync"
	"time"

	dspm_config "github.com/pgcrooks/dspm-scanner/internal/config"
	"github.com/pgcrooks/dspm-scanner/internal/finder"
)

type ScannerService struct {
	Instances int
}

type IScannerService interface {
	Run(ctx context.Context)
}

func InitScannerService(
	ctx context.Context,
	config *dspm_config.Config,
	bucketChan chan<- finder.BucketObjectBatch,
) (IScannerService, error) {
	slog.Info("init scanner service")

	service := ScannerService{
		Instances: config.Scanner.Instances,
	}

	return service, nil
}

func runScanner(ctx context.Context, id int) {
	slog.Info("run scanner instance", "id", id)

	run := true
	for run {
		select {
		case <-ctx.Done():
			slog.Info("stopping scanner", "id", id)
			run = false

		default:
			//todo
		}

		time.Sleep(time.Second)
	}
	slog.Info("stopped scanner instance", "id", id)
}

func (ss ScannerService) Run(ctx context.Context) {
	slog.Info("running ScannerService")

	// Wait group for all the individual scanner
	wg := sync.WaitGroup{}

	// New context for scanners
	scannerCtx, scannerCancel := context.WithCancel(ctx)
	defer scannerCancel()

	// Handle stop signal
	go func() {
		<-ctx.Done()
		slog.Info("stopping ScannerService")
		scannerCancel()
	}()

	// Run the scanners
	for i := range ss.Instances {
		wg.Go(func() {
			runScanner(scannerCtx, i)
		})
	}

	wg.Wait()
	slog.Info("stopped ScannerService")
}
