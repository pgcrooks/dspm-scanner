package scanner_int

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"
)

type LocalAPI interface {
	ReadDir(path string) ([]os.DirEntry, error)
}

func ListLocalBucket(ctx context.Context, path string) (BucketObjectBatch, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("can not list bucket. %v", err)
	}

	var contents BucketObjectBatch
	for _, file := range files {
		// Ignore directories
		if !file.IsDir() {
			var fileSize int64 = 0
			fileInfo, err := os.Stat(path + file.Name())
			if err != nil {
				log.Printf("failed to get file stat for %s", file.Name())
			} else {
				fileSize = fileInfo.Size()
			}
			obj := BucketObject{
				Key:  file.Name(),
				Size: fileSize,
			}
			contents = append(contents, obj)
		}
	}
	return contents, nil
}

func RunScraperService(ctx context.Context, cfg Config, messageChan chan<- BucketObjectBatch) {
	slog.Info("starting ScraperService")

	run := true
	for run {
		select {
		case <-ctx.Done():
			slog.Info("stopping ScraperService")
			run = false

		default:
			contents, err := ListLocalBucket(ctx, cfg.Local.Path)
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
