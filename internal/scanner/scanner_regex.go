package scanner

import (
	"context"
	"log/slog"
	"regexp"
	"time"

	"github.com/pgcrooks/dspm-scanner/internal/finder"
)

type matcher struct {
	Name    string
	Pattern string
	Regex   *regexp.Regexp
}

type scannerRegex struct {
	Scanner
	Matchers []matcher
}

func getDefaultMatchers() []matcher {
	return []matcher{
		{Name: "password", Pattern: "(?i)password[:= ]"},
		{Name: "key", Pattern: "(?i)key[:= ]"},
		{Name: "token", Pattern: "(?i)token[:= ]"},
	}
}

func newScannerRegex(bucketChan <-chan finder.BucketObjectBatch) (IScanner, error) {
	scanner := scannerRegex{
		Scanner: Scanner{
			Name:     "regex",
			Mode:     Regex,
			DataChan: bucketChan,
		},
	}

	for _, matcher := range getDefaultMatchers() {
		r, err := regexp.Compile(matcher.Pattern)
		if err != nil {
			slog.Error("unable to compile regex", "err", err.Error())
		} else {
			matcher.Regex = r
			scanner.Matchers = append(scanner.Matchers, matcher)
		}
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
