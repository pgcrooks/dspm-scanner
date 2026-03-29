package datastore

import (
	"context"
	"log/slog"

	scanner_int "github.com/pgcrooks/dspm-scanner/internal"
)

func InitDSMemory(ctx context.Context) (IDataStore, error) {
	slog.Info("Creating memory DS")
	ds := DataStoreMemory{}
	return ds, nil
}

func (ds DataStoreMemory) Close() {
	// No-op
	slog.Info("ds memory close")
}

func (ds DataStoreMemory) Write(object scanner_int.BucketObject) {
	slog.Info("ds memory write", "key", object.Key, "size", object.Size)
}
