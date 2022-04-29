package database

import (
	"database/sql"
	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

//go:embed schema.sql
var schema string

func NewDB(path string) (*sql.DB, error) {
	var err error

	Db, err = sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if _, err := Db.Exec(schema); err != nil {
		return nil, err
	}

	return Db, nil
}
