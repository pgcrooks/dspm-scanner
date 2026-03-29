package scanner_int

import (
	"context"
	"database/sql"
)

type BucketObject struct {
	Key  string
	Size int64
}

type BucketObjectBatch []BucketObject

type Config struct {
	DataStore struct {
		Memory struct {
			Enabled bool
		}
		LocalDB struct {
			Enabled bool
			Path    string
		}
	}
	Scraper struct {
		Aws struct {
			Enabled    bool
			BucketName string
		}
		Ds struct {
			Driver string
			Path   string
		}
		Local struct {
			Enabled bool
			Path    string
		}
	}
}

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
