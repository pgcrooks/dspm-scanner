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

func newDSMemory(bucketChan <-chan finder.BucketObjectBatch) (dataStoreMemory, error) {
	slog.Info("creating memory DS")
	return dataStoreMemory{
		DataStore: DataStore{
			Name:       "memory",
			Type:       Memory,
			BucketChan: bucketChan,
		},
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
	slog.Info("write", "key", object.Key, "size", object.Size, "p", fmt.Sprintf("%p", &dsm))
	dsm.Records = append(dsm.Records, record{
		Key:    object.Key,
		Size:   object.Size,
		Source: finder.LocalFS,
		Status: Idle,
		Result: None,
	})
}

func (dsm *dataStoreMemory) Stats() string {
	return fmt.Sprintf(
		"dataStoreMemory: len=%d cap=%d p=%p",
		len(dsm.Records),
		cap(dsm.Records),
		&dsm,
	)
}
