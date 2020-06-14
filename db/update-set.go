package db

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"github.com/sn-edit/sn-edit/conf"
)

func WriteUpdateSet(updateSetName string, updateSetSysID string, updateSetScope int64, current bool) error {
	dbc := conf.GetDB()
	// check if entry exists
	if exists, _ := UpdateSetExists(updateSetSysID); exists == true {
		log.WithFields(log.Fields{"error": "scope_exists"}).Debug("Update set already exists, no insert!")
		return nil
	}

	stmt, err := dbc.Prepare("INSERT INTO update_set(sys_id, name, sys_scope, current) VALUES(?,?,?,?)")
	defer stmt.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while preparing the query!")
		return err
	}

	_, err = stmt.Exec(updateSetSysID, updateSetName, updateSetScope, current)

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

func ListUpdateSets(scopeID int64) ([]map[string]interface{}, error) {
	dbc := conf.GetDB()

	rows, err := dbc.Query("SELECT sys_id, name, current FROM update_set WHERE sys_scope=?", scopeID)
	defer rows.Close()

	if err != nil {
		if err != sql.ErrNoRows {
			log.WithFields(log.Fields{"error": err}).Error("There was an error with the query!")
			return nil, err
		}
	}

	var updateSets []map[string]interface{}

	for rows.Next() {
		sysID := ""
		name := ""
		current := false

		err := rows.Scan(&sysID, &name, &current)

		updateSets = append(updateSets, map[string]interface{}{"sys_id": sysID, "name": name, "current": current})

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("There was an error while iterating through the results!")
			return nil, err
		}
	}

	return updateSets, nil
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

func UpdateSetsLoaded(scopeID int64) (bool, error) {
	dbc := conf.GetDB()
	stmt, err := dbc.Prepare("SELECT id FROM update_set WHERE sys_scope=? LIMIT 1")
	defer stmt.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
		return false, err
	}

	id := ""
	err = stmt.QueryRow(scopeID).Scan(&id)

	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Debug("The update set was not found in the database!")
		if err == sql.ErrNoRows {
			// no rows found, it does not exist
			return false, err
		} else {
			log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
			return false, err
		}
	}

	return true, nil
}
