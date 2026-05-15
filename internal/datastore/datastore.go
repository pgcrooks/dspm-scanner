package datastore

import (
	"context"
	"fmt"
	"log/slog"

	dspm_config "github.com/pgcrooks/dspm-scanner/internal/config"
	finder "github.com/pgcrooks/dspm-scanner/internal/finder"
)

type DataStoreType int

const (
	LocalDB DataStoreType = iota
	Memory
)

var dataStoreName = map[DataStoreType]string{
	LocalDB: "localdb",
	Memory:  "memory",
}

func (dst DataStoreType) String() string {
	return dataStoreName[dst]
}

// DataStore base class
type DataStore struct {
	Name       string
	Type       DataStoreType
	BucketChan <-chan finder.BucketObjectBatch
}

type IDataStore interface {
	Run(ctx context.Context)
	Close()
}

func NewDataStore(
	ctx context.Context,
	config *dspm_config.Config,
	messageChan <-chan finder.BucketObjectBatch,
) (IDataStore, error) {
	slog.Info("init datastore")

	// Config validator will ensure only one DS is enabled
	if config.DataStore.LocalDB.Enabled {
		ds, err := newDSLocalDB(config.DataStore.LocalDB.Path, messageChan)
		return &ds, err
	} else if config.DataStore.Memory.Enabled {
		ds, err := newDSMemory(messageChan)
		return &ds, err
	}

	return nil, fmt.Errorf("ds not implemented")
}
