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
		conf.Err("Error while preparing the query!", log.Fields{"error": err}, false)
		return err
	}

	_, err = stmt.Exec(updateSetSysID, updateSetName, updateSetScope, current)

	if err != nil {
		conf.Err("Error while executing the query!", log.Fields{"error": err}, false)
		return err
	}

	return nil
}

func QueryUpdateSet(updateSetSysID string) (bool, string, string) {
	dbc := conf.GetDB()
	stmt, err := dbc.Prepare("SELECT sys_id,name FROM update_set WHERE sys_id=? LIMIT 1")
	defer stmt.Close()

	if err != nil {
		conf.Err("Error while querying database data!", log.Fields{"error": err}, false)
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
			conf.Err("Error while querying database data!", log.Fields{"error": err}, false)
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
			conf.Err("Error while querying database data!", log.Fields{"error": err}, false)
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
			conf.Err("Error while iterating through the results!", log.Fields{"error": err}, false)
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
		conf.Err("Error while querying database data!", log.Fields{"error": err}, false)
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
			conf.Err("Error while querying database data!", log.Fields{"error": err}, false)
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
		conf.Err("Error while querying database data!", log.Fields{"error": err}, false)
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
			conf.Err("Error while querying database data!", log.Fields{"error": err}, false)
			return false, err
		}
	}

	return true, nil
}
