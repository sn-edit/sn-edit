package db

import (
	"database/sql"
	"errors"
	"github.com/0x111/sn-edit/conf"
	"github.com/0x111/sn-edit/file"
	log "github.com/sirupsen/logrus"
	"time"
)

// provides methods to handle entries
func WriteEntry(tableName string, sysName string, sysID string, sysScopeName string) error {
	dbc := conf.GetDB()
	// get table id from name if found
	// write table data
	err := WriteTable(tableName)

	if err != nil {
		log.WithFields(log.Fields{"warn": "table_write_error"}).Debug("Table already exists, no insert!")
	}

	success, tableID := QueryTable(tableName)

	if !success {
		err = errors.New("table_not_found")
		log.WithFields(log.Fields{"err": err}).Debug("Table not found! Please re-download!")
		return err
	}

	// write scope for file
	sysScope, err := RequestScopeDataFromInstance(sysScopeName)

	if err != nil {
		return err
	}

	// get scope for file
	found, fileScope := QueryScope(sysScope)

	if !found {
		err = errors.New("entry_not_found")
		log.WithFields(log.Fields{"err": err}).Debug("Entry not found! Please re-download!")
		return err
	}

	// check if entry exists
	if exists := EntryExists(tableID, sysID, fileScope); exists == true {
		log.WithFields(log.Fields{"warn": "entry_write_error"}).Debug("Entry already exists, no insert!")
		return nil
	}

	// filter name before entry to the db
	sysName = file.FilterSpecialChars(sysName)

	stmt, err := dbc.Prepare("INSERT INTO entry(sys_id, sys_name, entry_table, sys_scope, last_modified) VALUES(?,?,?,?,?)")
	defer stmt.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while preparing the query!")
		return err
	}

	_, err = stmt.Exec(sysID, sysName, tableID, fileScope, time.Now().UnixNano())

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while executing the query!")
		return err
	}

	return nil
}

func QuerySysName(tableName string, sysID string) (bool, string) {
	dbc := conf.GetDB()
	stmt, err := dbc.Prepare("SELECT sys_name FROM entry e LEFT JOIN entry_table t ON e.entry_table=t.id WHERE e.sys_id=? AND t.name=? LIMIT 1")
	defer stmt.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
		return false, ""
	}

	sysName := ""
	err = stmt.QueryRow(sysID, tableName).Scan(&sysName)

	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Debug("The script was not found in the database!")
		if err == sql.ErrNoRows {
			// no rows found, it does not exist
			return false, ""
		} else {
			log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
			return false, ""
		}
	}

	return true, sysName
}

func GetEntryScopeName(tableName string, sysID string) (bool, string) {
	dbc := conf.GetDB()
	stmt, err := dbc.Prepare("SELECT s.name FROM entry e LEFT JOIN entry_scope s ON e.sys_scope=s.id LEFT JOIN entry_table t ON e.entry_table=t.id WHERE e.sys_id=? AND t.name=? LIMIT 1;")
	defer stmt.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
		return false, ""
	}

	scopeName := ""
	err = stmt.QueryRow(sysID, tableName).Scan(&scopeName)

	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Debug("The script was not found in the database!")
		if err == sql.ErrNoRows {
			// no rows found, it does not exist
			return false, ""
		} else {
			log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
			return false, ""
		}
	}

	return true, scopeName
}

func EntryExists(tableID string, sysID string, sysScope string) bool {
	dbc := conf.GetDB()
	stmt, err := dbc.Prepare("SELECT sys_name FROM entry WHERE sys_id=? AND entry_table=? AND sys_scope=? LIMIT 1")
	defer stmt.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
		return false
	}

	id := ""
	err = stmt.QueryRow(sysID, tableID, sysScope).Scan(&id)

	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Debug("The script was not found in the database!")
		if err == sql.ErrNoRows {
			// no rows found, it does not exist
			return false
		} else {
			log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
			return false
		}
	}

	return true
}
