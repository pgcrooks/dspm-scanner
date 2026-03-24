package scanner_int

import "database/sql"

type BucketObject struct {
	Key  string
	Size int64
}

type BucketObjectBatch []BucketObject

type Config struct {
	Aws struct {
		Enabled    bool
		BucketName string
	}
	Db struct {
		Driver string
		Path   string
	}
	Local struct {
		Enabled bool
		Path    string
	}
}

type DataStoreType int

const (
	LocalDB DataStoreType = iota
)

type DataStore struct {
	Type    DataStoreType
	LocalDB *sql.DB
}
