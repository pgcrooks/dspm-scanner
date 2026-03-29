package datastore

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	scanner_int "github.com/pgcrooks/dspm-scanner/internal"
	finder "github.com/pgcrooks/dspm-scanner/internal/finder"
)

type DataStoreType int

const (
	LocalDB DataStoreType = iota
	Memory
)

type DataStore struct {
	Name string
	Type DataStoreType
}

type DataStoreLocalDB struct {
	DataStore
	SqlDB *sql.DB
}

type DataStoreMemory struct {
	DataStore
	// Actual store lives here
}

type IDataStore interface {
	GetName() string
	Close()
	Write(object finder.BucketObject)
}

var dataStoreName = map[DataStoreType]string{
	LocalDB: "localdb",
	Memory:  "memory",
}

func (dst DataStoreType) String() string {
	return dataStoreName[dst]
}

func (ds DataStore) GetName() string {
	return ds.Name
}

func (ds DataStore) Write(object finder.BucketObject) {
	slog.Error("ds write not impl")
}

func InitDataStore(ctx context.Context, config *scanner_int.Config) (IDataStore, error) {
	// Config validator will ensure only one DS is enabled
	if config.DataStore.LocalDB.Enabled {
		ds, err := initDSLocalDB(ctx, config.DataStore.LocalDB.Path)
		return ds, err
	} else if config.DataStore.Memory.Enabled {
		ds, err := InitDSMemory(ctx)
		return ds, err
	} else {
		return DataStore{}, fmt.Errorf("ds not implemented")
	}
}

func (ds DataStore) Close() {
	// Can be overridden by child impls
	slog.Info("base ds close")
}

func RunDataService(ctx context.Context, ids IDataStore, messageChan <-chan finder.BucketObjectBatch) {
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
				ids.Write(object)
			}
		}

		time.Sleep(time.Second)
	}

	slog.Info("terminated DataService")
}
