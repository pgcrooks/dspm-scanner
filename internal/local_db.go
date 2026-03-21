package scanner_int

import (
	"database/sql"
	"fmt"

	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

const ARTIFACT_TABLE string = "artifacts"

func InitLocalDB(dbName string) (*sql.DB, error) {
	slog.Info("Opening", "dbName", dbName)
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	// Create table
	sqlStmt := fmt.Sprintf(
		"create table if not exists %s (id integer not null primary key, key text not null, size int64, scanned bool default FALSE);",
		ARTIFACT_TABLE,
	)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("can not create table. %v", err)
	}
	return db, nil
}

func CloseLocalDB(db *sql.DB) {
	db.Close()
}
