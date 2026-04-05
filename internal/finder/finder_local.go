package finder

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"
)

type finderLocal struct {
	Finder
	Path string
}

func newFinderLocal(path string) (IFinder, error) {
	return &finderLocal{
		Finder: Finder{
			Name: "local",
		},
		Path: path,
	}, nil
}

func (f finderLocal) Run(ctx context.Context) {
	slog.Info("running finder local", "obj", f)

	run := true
	for run {
		select {
		case <-ctx.Done():
			slog.Info("stopping finder local", "obj", f)
			run = false

		default:
			files, err := os.ReadDir(f.Path)
			if err == nil {
				for _, file := range files {
					// Ignore directories
					if !file.IsDir() {
						var fileSize int64 = 0
						fileInfo, err := os.Stat(f.Path + file.Name())
						if err != nil {
							log.Printf("failed to get file stat for %s", file.Name())
						} else {
							fileSize = fileInfo.Size()
						}
						obj := BucketObject{
							Key:  file.Name(),
							Size: fileSize,
						}
						slog.Info("found metadata", "file", obj)
						//todo write to channel
					}
				}
			} else {
				slog.Warn("cannot read path", "err", err.Error())
			}
		}

		time.Sleep(time.Second)
	}
	slog.Info("stopped finder local", "obj", f)
}
