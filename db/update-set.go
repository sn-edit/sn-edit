package db

import (
	"database/sql"
	"github.com/0x111/sn-edit/conf"
	log "github.com/sirupsen/logrus"
)

func WriteUpdateSet(updateSetName string, updateSetSysID string, updateSetScope string) error {
	dbc := conf.GetDB()
	// check if entry exists
	if exists, _ := UpdateSetExists(updateSetSysID); exists == true {
		log.WithFields(log.Fields{"error": "scope_exists"}).Debug("Update set already exists, no insert!")
		return nil
	}

	stmt, err := dbc.Prepare("INSERT INTO update_set(sys_id, name, sys_scope) VALUES(?,?,?)")
	defer stmt.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while preparing the query!")
		return err
	}

	_, err = stmt.Exec(updateSetSysID, updateSetName, updateSetScope)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while executing the query!")
		return err
	}

	return nil
}

func QueryUpdateSet(updateSetSysID string) (bool, string, string) {
	dbc := conf.GetDB()
	stmt, err := dbc.Prepare("SELECT sys_id,name FROM update_set WHERE sys_id=? LIMIT 1")
	defer stmt.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
		return false, "", ""
	}

	sysID := ""
	name := ""
	err = stmt.QueryRow(updateSetSysID).Scan(&sysID, &name)

	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Debug("The scope was not found in the database!")
		if err == sql.ErrNoRows {
			// no rows found, it does not exist
			return false, "", ""
		} else {
			log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
			return false, "", ""
		}
	}

	return true, sysID, name
}

func UpdateSetExists(updateSetSysID string) (bool, string) {
	dbc := conf.GetDB()
	stmt, err := dbc.Prepare("SELECT id FROM update_set WHERE sys_id=? LIMIT 1")
	defer stmt.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
		return false, ""
	}

	id := ""
	err = stmt.QueryRow(updateSetSysID).Scan(&id)

	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Debug("The scope was not found in the database!")
		if err == sql.ErrNoRows {
			// no rows found, it does not exist
			return false, ""
		} else {
			log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
			return false, ""
		}
	}

	return true, id
}
