package finder

import (
	"context"
	"log/slog"
	"time"

	scanner_int "github.com/pgcrooks/dspm-scanner/internal"
)

type BucketObject struct {
	Key  string
	Size int64
}

type BucketObjectBatch []BucketObject

type IFinderService interface {
	Run()
}

type IFinderImpl interface {
	Run()
}

type FinderLocal struct {
}

type FinderService struct {
	Local FinderLocal
}

// func InitFinderService(ctx context.Context, config *scanner_int.Config) (IFinderService, error) {
// 	slog.Info("init finder service")
// 	finderService := IFinderService{}
// 	if config.Scraper.Local.Enabled {
// 		slog.Info("finder local enabled")
// 	}
// 	return
// }

func RunFinderService(ctx context.Context, cfg scanner_int.Config, messageChan chan<- BucketObjectBatch) {
	slog.Info("starting FinderService")

	run := true
	for run {
		select {
		case <-ctx.Done():
			slog.Info("stopping FinderService")
			run = false

		default:
			contents, err := listLocalBucket(ctx, cfg.Scraper.Local.Path)
			if err != nil {
				slog.Error(err.Error())
			} else {
				messageChan <- contents
			}
		}

		time.Sleep(time.Second)
	}

	slog.Info("terminated ScraperService")
}
