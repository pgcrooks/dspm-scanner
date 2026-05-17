package datastore

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	finder "github.com/pgcrooks/dspm-scanner/internal/finder"
)

type dataStoreMemory struct {
	DataStore
	Records map[string]ObjectRecord
	LastID  int
}

func newDSMemory(bucketChan <-chan finder.BucketObjectBatch) (dataStoreMemory, error) {
	slog.Info("creating memory DS")
	return dataStoreMemory{
		DataStore: DataStore{
			Name:       "memory",
			Type:       Memory,
			BucketChan: bucketChan,
		},
		Records: make(map[string]ObjectRecord),
		LastID:  0,
	}, nil
}

func (dsm *dataStoreMemory) Close() {
	// No-op
	slog.Info("closed")
}

func (dsm *dataStoreMemory) Run(ctx context.Context) {
	slog.Info("running dataStoreMemory", "obj", dsm)

	run := true
	for run {
		select {
		case <-ctx.Done():
			slog.Info("stopping dataStoreMemory")
			run = false

		case msg1 := <-dsm.BucketChan:
			for _, object := range msg1 {
				slog.Info("rx", "key", object.Key, "size", object.Size)
				dsm.Write(object)
			}
			// Be loud for now
			slog.Info(dsm.Stats())

		default:
			//TODO run separate read and write coroutines
		}

		time.Sleep(time.Second)
	}
}

func (dsm *dataStoreMemory) Write(object finder.BucketObject) {
	slog.Info("write", "key", object.Key, "size", object.Size)

	// Check if path already exists, and update if necessary
	// A cache should really be used for this
	for recordID, record := range dsm.Records {
		if object.Key == record.Key {
			if object.Size == record.Size {
				// Skip, already exists and is up to date
				slog.Info("already up to date", "recordID", recordID)
				return
			} else {
				// Exists but size has changed, delete and recreate
				slog.Info("refreshing", "recordID", recordID)
				delete(dsm.Records, recordID)

				// No need to check the rest
				break
			}
		}
	}

	id := dsm.makeID()
	slog.Info("creating", "recordID", id)
	dsm.Records[id] = ObjectRecord{
		ID:     id,
		Hash:   "",
		Key:    object.Key,
		Size:   object.Size,
		Source: finder.LocalFS,
	}
}

func (dsm *dataStoreMemory) Stats() string {
	return fmt.Sprintf("dataStoreMemory: len=%d", len(dsm.Records))
}

func (dsm *dataStoreMemory) makeID() string {
	// TODO worry about wrapping
	dsm.LastID++
	return fmt.Sprintf("ds-%08d", dsm.LastID)
}
