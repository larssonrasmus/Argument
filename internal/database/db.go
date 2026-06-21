package database

import (
	"database/sql"

	"argument/internal/config"

	// Imported to get MySQL support
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	Connect(config.Conf.Database.URI)
}

// Connect prepares the database connection
func Connect(URI string) {
	if db != nil {
		return
	}

	var err error
	db, err = sql.Open("mysql", URI)

	if err == nil {
		// Check that we actually managed to get a connection
		err = db.Ping()
	}

	if err != nil {
		panic(err)
	}
}

// Disconnect the database
func Disconnect() {
	if db != nil {
		db.Close()
		db = nil
	}
}
