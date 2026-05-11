package datastore

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	finder "github.com/pgcrooks/dspm-scanner/internal/finder"
)

type recordStatus int

const (
	Idle recordStatus = iota
	Scanning
	Scanned
)

type recordResult int

const (
	None recordResult = iota
	Clean
	Alert
)

type record struct {
	Key    string
	Size   int64
	Source finder.FinderType
	Status recordStatus
	Result recordResult
}

type dataStoreMemory struct {
	DataStore
	Type    DataStoreType
	Records []record
}

func InitDSMemory(ctx context.Context, bucketChan <-chan finder.BucketObjectBatch) (IDataStore, error) {
	slog.Info("Creating memory DS")
	// ds := DataStoreMemory{}
	// return ds, nil
	return &dataStoreMemory{
		DataStore: DataStore{
			Name:       "memory",
			BucketChan: bucketChan,
		},
	}, nil
}

func (ds dataStoreMemory) Close() {
	// No-op
	slog.Info("ds memory close")
}

func (ds dataStoreMemory) Run(ctx context.Context) {
	slog.Info("running DataStoreMemory", "obj", ds)

	run := true
	for run {
		select {
		case <-ctx.Done():
			slog.Info("stopping DataStoreMemory")
			run = false

		case msg1 := <-ds.BucketChan:
			for _, object := range msg1 {
				slog.Info("rx", "key", object.Key, "size", object.Size)
				ds.Write(object)
			}
			slog.Info(ds.Stats())

		default:
			//TODO run separate read and write coroutines
		}

		time.Sleep(time.Second)
	}
}

func (ds dataStoreMemory) Write(object finder.BucketObject) {
	slog.Info("write", "key", object.Key, "size", object.Size)
	ds.Records = append(ds.Records, record{
		Key:    object.Key,
		Size:   object.Size,
		Source: finder.LocalFS,
		Status: Idle,
		Result: None,
	})
}

func (ds dataStoreMemory) Stats() string {
	return fmt.Sprintf(
		"DataStoreMemory: len=%d cap=%d",
		len(ds.Records),
		cap(ds.Records),
	)
}
