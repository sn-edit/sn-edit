package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/0x111/sn-edit/api"
	"github.com/0x111/sn-edit/conf"
	"github.com/icza/dyno"
	log "github.com/sirupsen/logrus"
)

func WriteScope(sysID string, scopeApiURL string) error {
	//config := conf.GetConfig()
	dbc := conf.GetDB()

	// get the table details from REST
	// setup the table API URL url
	response, err := api.Get(scopeApiURL)

	if err != nil {
		return err
	}

	// unmarshal response
	var responseResult map[string]interface{}
	err = json.Unmarshal(response, &responseResult)

	if err != nil {
		fmt.Println("There was an error while unmarshalling the response!", err)
		return err
	}

	result, err := dyno.Get(responseResult, "result")

	if err != nil {
		fmt.Println("Error getting the result key!", err)
		return err
	}

	// get scope name
	tableScopeName, err := dyno.GetString(result, "scope")

	// check if entry exists
	if exists, _ := ScopeExists(tableScopeName); exists == true {
		log.WithFields(log.Fields{"error": "scope_exists"}).Debug("Scope already exists, no insert!")
		return nil
	}

	stmt, err := dbc.Prepare("INSERT INTO entry_scope(sys_id, name) VALUES(?,?)")
	defer stmt.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while preparing the query!")
		return err
	}

	_, err = stmt.Exec(sysID, tableScopeName)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while executing the query!")
		return err
	}

	return nil
}

func QueryScope(sysID string) (bool, string) {
	dbc := conf.GetDB()
	stmt, err := dbc.Prepare("SELECT id FROM entry_scope WHERE sys_id=? LIMIT 1")
	defer stmt.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
		return false, ""
	}

	id := ""
	err = stmt.QueryRow(sysID).Scan(&id)

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

func ScopeExists(scopeName string) (bool, string) {
	dbc := conf.GetDB()
	stmt, err := dbc.Prepare("SELECT id FROM entry_scope WHERE name=? LIMIT 1")
	defer stmt.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
		return false, ""
	}

	sysID := ""
	err = stmt.QueryRow(scopeName).Scan(&sysID)

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

	return true, sysID
}

func GetScopeNameFromSysID(sysID string) (bool, string) {
	dbc := conf.GetDB()
	stmt, err := dbc.Prepare("SELECT name FROM entry_scope WHERE sys_id=? LIMIT 1")
	defer stmt.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
		return false, ""
	}

	name := ""
	err = stmt.QueryRow(sysID).Scan(&name)

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

	return true, name
}
