package scanner_int

import (
	"context"
	"fmt"
)

type DataStoreAPI interface {
	Close()
}

func InitDataStore(ctx context.Context, dsType DataStoreType) (DataStore, error) {
	switch dsType {
	case LocalDB:
		// TODO: hardcoded until Config module exists
		localDB, err := InitLocalDB("dspm.db")
		if err != nil {
			return DataStore{}, err
		}
		var ds DataStore
		ds.Type = LocalDB
		ds.LocalDB = localDB
		return ds, nil
	default:
		return DataStore{}, fmt.Errorf("unknown ds: %d", dsType)
	}
}

func (ds DataStore) Close() {
	switch ds.Type {
	case LocalDB:
		CloseLocalDB(ds.LocalDB)
	default:
		// No-op
	}
}
