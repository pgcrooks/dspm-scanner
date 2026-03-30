package datastore

import (
	"context"
	"database/sql"
	"fmt"

	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

const ARTIFACT_TABLE string = "artifacts"

func initDSLocalDB(ctx context.Context, dbName string) (IDataStore, error) {
	slog.Info("Opening", "dbName", dbName)
	ds := DataStoreLocalDB{}
	ds.Name = "fooDB"
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return &ds, err
	}

	// Create table
	sqlStmt := fmt.Sprintf(
		"create table if not exists %s (id integer not null primary key, key text not null, size int64, scanned bool default FALSE);",
		ARTIFACT_TABLE,
	)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return ds, fmt.Errorf("can not create table. %v", err)
	}

	ds.SqlDB = db
	ds.Type = LocalDB
	return ds, nil
}

func (ds DataStoreLocalDB) Close() {
	err := ds.SqlDB.Close()
	if err != nil {
		slog.Error("unable to close db", "err", err)
	}
}
