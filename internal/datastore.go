package scanner_int

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
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
	Run(ctx context.Context, messageChan <-chan BucketObjectBatch)
}

type DataStoreAPI interface {
	Close()
	RunDataService()
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

func InitDataStore(ctx context.Context, config *Config) (IDataStore, error) {
	// Config validator will ensure only one DS is enabled
	if config.DataStore.LocalDB.Enabled {
		ds, err := InitDSLocalDB(config.DataStore.LocalDB.Path)
		return ds, err
	} else if config.DataStore.Memory.Enabled {
		ds, err := InitDSMemory()
		return ds, err
	} else {
		return DataStore{}, fmt.Errorf("ds not implemented")
	}
}

func (ds DataStore) Close() {
	// Useless wrapper for now
	ds.Close()
}

func (ds DataStore) Run(ctx context.Context, messageChan <-chan BucketObjectBatch) {
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
