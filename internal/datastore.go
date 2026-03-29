package scanner_int

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type DataStoreAPI interface {
	Close()
	RunDataService()
}

func InitDataStore(ctx context.Context, config *Config) (DataStore, error) {
	switch config.Ds.Driver {
	case "sqlite":
		db, err := InitLocalDB(config.Ds.Path)
		if err != nil {
			return DataStore{}, err
		}
		ds := DataStore{}
		ds.Type = LocalDB
		ds.LocalDB = db
		return ds, nil
	default:
		return DataStore{}, fmt.Errorf("unknown ds: %s", config.Ds.Driver)
	}
}

func (ds DataStore) Close() {
	switch ds.Type {
	case LocalDB:
		CloseLocalDB(ds.LocalDB)
	default:
		// No-op
	}
}

func (ds DataStore) RunDataService(ctx context.Context, messageChan <-chan BucketObjectBatch) {
	slog.Info("starting DataService")

	run := true
	for run {
		select {
		case <-ctx.Done():
			slog.Info("stopping DataService")
			run = false

		case msg1 := <-messageChan:
			for _, object := range msg1 {
				slog.Info("rx scrape_data", "key", object.Key, "size", object.Size)
			}
		}

		time.Sleep(time.Second)
	}

	slog.Info("terminated DataService")
}
