package conf

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var (
	db *sql.DB
)

func GetDB() *sql.DB {
	return db
}

func ConnectDB() *sql.DB {
	var err error
	config := GetConfig()

	log.Debug("Connecting to the database...")
	path := config.GetString("app.core.db.path")

	if path == "" {
		Err("Please specify a valid database path!", log.Fields{"error": err, "path": path}, true)
	}

	db, err = sql.Open("sqlite3", "file:"+path+"?cache=shared")

	// Check error for database connection
	if err != nil {
		Err("Please specify a valid database path!", log.Fields{"error": err}, true)
	}

	return db
}
