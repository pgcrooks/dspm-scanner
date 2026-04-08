package finder

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	dspm_config "github.com/pgcrooks/dspm-scanner/internal/config"
)

// Individual data object's metadata
type BucketObject struct {
	Key  string
	Size int64
}

// Batch of data object metadata
type BucketObjectBatch []BucketObject

type FinderType int

const (
	LocalFS FinderType = iota
	AWSS3
)

var finderTypeName = map[FinderType]string{
	LocalFS: "local_fs",
	AWSS3:   "aws_s3",
}

func (ft FinderType) String() string {
	return finderTypeName[ft]
}

// Idividual Finder for a certain data source (Local, AWS S3...)
type IFinder interface {
	Run(ctx context.Context)
}

// Finder Service which owns many Finders
type FinderService struct {
	Finders []IFinder
}

// Finder Service interface
type IFinderService interface {
	Run(ctx context.Context)
}

// Finder base class
type Finder struct {
	Name       string
	BucketChan chan<- BucketObjectBatch
}

func InitFinderService(
	ctx context.Context,
	config *dspm_config.Config,
	bucketChan chan<- BucketObjectBatch,
) (IFinderService, error) {
	slog.Info("init finder service")

	// Error checking
	if !config.Finder.Aws.Enabled && !config.Finder.Local.Enabled {
		return nil, fmt.Errorf("no finders enabled")
	}

	// Spin up each finder
	service := FinderService{}

	if config.Finder.Aws.Enabled {
		slog.Info("finder enabled: aws")
		finder, err := newFinderAWSS3(ctx, config.Finder.Aws.BucketName, bucketChan)
		if err != nil {
			slog.Warn("cannot create aws s3 finder", "err", err.Error())
		} else {
			service.Finders = append(service.Finders, finder)
		}
	}

	if config.Finder.Local.Enabled {
		slog.Info("finder enabled: local")
		finder, err := newFinderLocal(config.Finder.Local.Path, bucketChan)
		if err != nil {
			slog.Warn("cannot create local finder", "err", err.Error())
		} else {
			service.Finders = append(service.Finders, finder)
		}
	}

	slog.Info("finder service initialised", "numFinders", len(service.Finders))
	return service, nil
}

func (fs FinderService) Run(ctx context.Context) {
	slog.Info("running FinderService")

	// Wait group for all the individual finders
	wg := sync.WaitGroup{}

	// New context for finders
	finderCtx, finderCancel := context.WithCancel(ctx)
	defer finderCancel()

	// Handle stop signal
	go func() {
		<-ctx.Done()
		slog.Info("stopping FinderService")
		finderCancel()
	}()

	// Run the finders
	for _, f := range fs.Finders {
		wg.Go(func() {
			f.Run(finderCtx)
		})
	}

	wg.Wait()
	slog.Info("stopped FinderService")
}
