package model

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func OpenDB(dsn string) error {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return err
	}
	cfg.Loc = time.Local
	cfg.ParseTime = true
	dsn = cfg.FormatDSN()

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return nil
}

func CloseDB() error {
	return db.Close()
}
