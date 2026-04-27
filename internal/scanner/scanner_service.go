package scanner

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	dspm_config "github.com/pgcrooks/dspm-scanner/internal/config"
	"github.com/pgcrooks/dspm-scanner/internal/finder"
)

type ScannerType int

const (
	Regex ScannerType = iota
)

var scannerTypeName = map[ScannerType]string{
	Regex: "regex",
}

func (st ScannerType) String() string {
	return scannerTypeName[st]
}

// Scanner Service owns many Scanners
type ScannerService struct {
	Instances int
	Scanners  []IScanner
}

type IScannerService interface {
	Run(ctx context.Context)
}

// Scanner base class
type Scanner struct {
	Name     string
	Mode     ScannerType
	DataChan <-chan finder.BucketObjectBatch
}

type IScanner interface {
	Run(ctx context.Context)
}

func InitScannerService(
	ctx context.Context,
	config *dspm_config.Config,
	bucketChan <-chan finder.BucketObjectBatch,
) (IScannerService, error) {
	slog.Info("init scanner service")

	// Error checking
	if !config.Scanner.Regex.Enabled {
		return nil, fmt.Errorf("no scanners enabled")
	}

	// Spin up each scanner
	service := ScannerService{
		Instances: config.Scanner.Instances,
	}

	if config.Scanner.Regex.Enabled {
		slog.Info("scanner enabled: regex")
		scanner, err := newScannerRegex(bucketChan)
		if err != nil {
			slog.Warn("cannot create scanner regex", "err", err.Error())
		} else {
			service.Scanners = append(service.Scanners, scanner)
		}
	}

	return service, nil
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
	for _, s := range ss.Scanners {
		wg.Go(func() {
			s.Run(scannerCtx)
		})
	}

	wg.Wait()
	slog.Info("stopped ScannerService")
}
