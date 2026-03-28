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

func InitDataStore(ctx context.Context, dsType DataStoreType) (DataStore, error) {
	switch dsType {
	case LocalDB:
		// TODO: hardcoded until Config module exists
		localDB, err := InitLocalDB("dspm.db")
		if err != nil {
			return DataStore{}, err
		}
		var ds DataStore
		ds.Type = LocalDB
		ds.LocalDB = localDB
		return ds, nil
	default:
		return DataStore{}, fmt.Errorf("unknown ds: %d", dsType)
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
