package conf

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"os"
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
	path := config.GetString("app.db.path")

	if path == "" {
		log.WithFields(log.Fields{"error": err, "path": path}).Error("Please specify a database path!")
	}

	db, err = sql.Open("sqlite3", "file:"+path+"?cache=shared")

	// Check error for database connection
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was a problem while connecting to the database!")
		os.Exit(1)
	}

	return db
}
