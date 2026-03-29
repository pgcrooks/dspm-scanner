package datastore

import "log/slog"

func InitDSMemory() (IDataStore, error) {
	slog.Info("Creating memory DS")
	ds := DataStoreMemory{}
	return ds, nil
}

func (ds DataStoreMemory) Close() {
	// No-op
}
