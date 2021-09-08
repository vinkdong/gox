package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"v2k.io/gox/log"
)

type Database struct {
	DataSourceName string
}

func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	db, err := sql.Open("mysql", d.DataSourceName)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error(err)
		}
	}()
	return db.Query(query, args...)
}

func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	db, err := sql.Open("mysql", d.DataSourceName)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error(err)
		}
	}()
	return db.Exec(query, args...)
}
