package datastore

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	finder "github.com/pgcrooks/dspm-scanner/internal/finder"
)

const ARTIFACT_TABLE string = "artifacts"

type dataStoreLocalDB struct {
	DataStore
	SqlDB *sql.DB
}

func newDSLocalDB(dbName string, bucketChan <-chan finder.BucketObjectBatch) (dataStoreLocalDB, error) {
	slog.Info("creating local db DS")
	slog.Info("Opening", "dbName", dbName)
	ds := dataStoreLocalDB{
		DataStore: DataStore{
			Name:       "sqlite",
			Type:       LocalDB,
			BucketChan: bucketChan,
		},
	}
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return ds, err
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
	return ds, nil
}

func (dsl *dataStoreLocalDB) Close() {
	err := dsl.SqlDB.Close()
	if err != nil {
		slog.Error("unable to close db", "err", err)
	}
}

func (dsl *dataStoreLocalDB) Run(ctx context.Context) {
	slog.Info("running dataStoreLocalDB", "obj", dsl)

	run := true
	for run {
		select {
		case <-ctx.Done():
			slog.Info("stopping dataStoreLocalDB")
			run = false

		case msg1 := <-dsl.BucketChan:
			for _, object := range msg1 {
				slog.Info("rx", "key", object.Key, "size", object.Size)
				dsl.Write(object)
			}

		default:
			//TODO run separate read and write coroutines
		}

		time.Sleep(time.Second)
	}
}

func (dsl *dataStoreLocalDB) Stats() string {
	return "not_impl"
}

func (dsl *dataStoreLocalDB) Write(object finder.BucketObject) {
	slog.Warn("not impl")
}
