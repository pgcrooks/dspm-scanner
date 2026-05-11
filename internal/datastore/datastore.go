package datastore

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

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
	BucketChan <-chan finder.BucketObjectBatch
}

// Individual DataStore object
type IDataStore interface {
	Close()
	Run(ctx context.Context)
	Stats() string
	Write(object finder.BucketObject)
}

// DataStore service
type DataStoreService struct {
	DSImpl IDataStore
}

// DataStore service interface
type IDataStoreService interface {
	Run(ctx context.Context)
}

func InitDataStoreService(
	ctx context.Context,
	config *dspm_config.Config,
	messageChan <-chan finder.BucketObjectBatch,
) (IDataStoreService, error) {
	slog.Info("init datastore service")

	service := DataStoreService{}

	// Config validator will ensure only one DS is enabled
	if config.DataStore.LocalDB.Enabled {
		ds, err := initDSLocalDB(ctx, config.DataStore.LocalDB.Path, messageChan)
		if err != nil {
			return nil, fmt.Errorf("cannot init localdb datastore")
		}
		service.DSImpl = ds
	} else if config.DataStore.Memory.Enabled {
		ds, err := InitDSMemory(ctx, messageChan)
		if err != nil {
			return nil, fmt.Errorf("cannot init memory datastore")
		}
		service.DSImpl = ds
	} else {
		return nil, fmt.Errorf("ds not implemented")
	}

	slog.Info("datastore service initialiserd")
	return service, nil
}

func (ds DataStore) Close() {
	// Can be overridden by child impls
	slog.Info("base ds close")
}

func (dss DataStoreService) Run(ctx context.Context) {
	slog.Info("running DataStoreService")

	wg := sync.WaitGroup{}

	dataStoreContext, dataStoreCancel := context.WithCancel(ctx)
	defer dataStoreCancel()

	// Handle stop signal
	go func() {
		<-ctx.Done()
		slog.Info("stopping DataStoreService")
		dataStoreCancel()
	}()

	// Run the datastore
	wg.Go(func() {
		dss.DSImpl.Run(dataStoreContext)
	})

	wg.Wait()

	slog.Info("terminated DataStoreService")
}
