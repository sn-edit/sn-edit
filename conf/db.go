package db

import (
	"database/sql"
	"github.com/0x111/sn-edit/conf"
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
	config := conf.GetConfig()

	log.Debug("Connecting to the database...")

	db, err = sql.Open("sqlite3", "file:"+config.GetString("app.db.path")+"?cache=shared")

	// Check error for database connection
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was a problem while connecting to the database!")
		os.Exit(1)
	}

	return db
}
