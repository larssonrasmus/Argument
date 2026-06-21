package database

import (
	"database/sql"

	"argument/internal/config"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func init() {
	Connect(config.Conf.Database.Driver, config.Conf.Database.URI)
}

// Connect prepares the database connection
func Connect(database string, URI string) {
	if db != nil {
		return
	}

	var err error
	db, err = sql.Open(database, URI)

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
