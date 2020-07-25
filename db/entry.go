package db

import (
	"database/sql"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/sn-edit/sn-edit/conf"
	"github.com/sn-edit/sn-edit/file"
	"time"
)

// provides methods to handle entries
func WriteEntry(tableName string, uniqueKeyName string, sysID string, sysScopeSysID string) error {
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
		log.WithFields(log.Fields{"warn": err}).Debug("Table not found! Please re-download!")
		return err
	}

	// write scope for file
	sysScope, err := RequestScopeDataFromInstance(sysScopeSysID)

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
	uniqueKeyName = file.FilterSpecialChars(uniqueKeyName)

	stmt, err := dbc.Prepare("INSERT INTO entry(sys_id, unique_key, entry_table, sys_scope, last_modified) VALUES(?,?,?,?,?)")
	defer stmt.Close()

	if err != nil {
		conf.Err("There was an error while preparing the query!", log.Fields{"error": err}, false)
		return err
	}

	_, err = stmt.Exec(sysID, uniqueKeyName, tableID, fileScope, time.Now().UnixNano())

	if err != nil {
		conf.Err("There was an error while executing the query!", log.Fields{"error": err}, false)
		return err
	}

	return nil
}

func QueryUniqueKey(tableName string, sysID string) (bool, string) {
	dbc := conf.GetDB()
	stmt, err := dbc.Prepare("SELECT unique_key FROM entry e LEFT JOIN entry_table t ON e.entry_table=t.id WHERE e.sys_id=? AND t.name=? LIMIT 1")
	defer stmt.Close()

	if err != nil {
		conf.Err("There was an error while querying the database!", log.Fields{"error": err}, false)
		return false, ""
	}

	uniqueKey := ""
	err = stmt.QueryRow(sysID, tableName).Scan(&uniqueKey)

	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Debug("The script was not found in the database!")
		if err == sql.ErrNoRows {
			// no rows found, it does not exist
			return false, ""
		} else {
			conf.Err("There was an error while querying the database!", log.Fields{"error": err}, false)
			return false, ""
		}
	}

	return true, uniqueKey
}

func GetEntryScopeName(tableName string, sysID string) (bool, string) {
	dbc := conf.GetDB()
	stmt, err := dbc.Prepare("SELECT s.name FROM entry e LEFT JOIN entry_scope s ON e.sys_scope=s.id LEFT JOIN entry_table t ON e.entry_table=t.id WHERE e.sys_id=? AND t.name=? LIMIT 1;")
	defer stmt.Close()

	if err != nil {
		conf.Err("There was an error while querying the database!", log.Fields{"error": err}, false)
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
			conf.Err("There was an error while querying the database!", log.Fields{"error": err}, false)
			return false, ""
		}
	}

	return true, scopeName
}

// todo: Implement update of existing entry with the updated fields
func EntryExists(tableID string, sysID string, sysScope int64) bool {
	dbc := conf.GetDB()
	stmt, err := dbc.Prepare("SELECT unique_key FROM entry WHERE sys_id=? AND entry_table=? AND sys_scope=? LIMIT 1")
	defer stmt.Close()

	if err != nil {
		conf.Err("There was an error while querying the database!", log.Fields{"error": err}, false)
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
			conf.Err("There was an error while querying the database!", log.Fields{"error": err}, false)
			return false
		}
	}

	return true
}
