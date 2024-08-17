package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DATABASE_URL = "krankentransport.db"

func GetDB() *sql.DB {

	db, err := sql.Open("sqlite3", DATABASE_URL)
	if err != nil {
		panic(err)
	}
	return db
}
