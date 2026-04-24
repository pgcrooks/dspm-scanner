package scanner

import (
	"context"
	"log/slog"
	"regexp"
	"time"

	"github.com/pgcrooks/dspm-scanner/internal/finder"
)

type scannerRegex struct {
	Scanner
	Checks []*regexp.Regexp
}

func newScannerRegex(bucketChan <-chan finder.BucketObjectBatch) (IScanner, error) {
	scanner := scannerRegex{
		Scanner: Scanner{
			Name:     "regex",
			Mode:     Regex,
			DataChan: bucketChan,
		},
	}

	r, err := regexp.Compile("(?i)password[:= ]")
	if err != nil {
		slog.Error("unable to compile regex", "err", err.Error())
	} else {
		scanner.Checks = append(scanner.Checks, r)
	}

	return &scanner, nil
}

func (s scannerRegex) Run(ctx context.Context) {
	slog.Info("running scanner regex", "obj", s)

	run := true
	for run {
		select {
		case <-ctx.Done():
			slog.Info("stopping scanner regex", "obj", s)
			run = false

		case dataMessage := <-s.DataChan:
			for _, object := range dataMessage {
				slog.Info("rx data into scanner regex", "key", object.Key)
			}

		default:
			// Should not get here
		}

		time.Sleep(time.Second)
	}
	slog.Info("stopped scanner regex", "obj", s)
}
