package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/0x111/sn-edit/api"
	"github.com/0x111/sn-edit/conf"
	"github.com/icza/dyno"
	log "github.com/sirupsen/logrus"
	"strings"
)

func WriteScope(sysID string, scopeName string) error {
	//config := conf.GetConfig()
	dbc := conf.GetDB()

	stmt, err := dbc.Prepare("INSERT INTO entry_scope(sys_id, name) VALUES(?,?)")
	defer stmt.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while preparing the query!")
		return err
	}

	_, err = stmt.Exec(sysID, scopeName)

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

func ScopeExists(scopeName string) (bool, string, string) {
	dbc := conf.GetDB()
	stmt, err := dbc.Prepare("SELECT id,sys_id FROM entry_scope WHERE name=? LIMIT 1")
	defer stmt.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
		return false, "", ""
	}

	id := ""
	sysID := ""
	err = stmt.QueryRow(scopeName).Scan(&id, &sysID)

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

	return true, id, sysID
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

func RequestScopeData(scopeApiURL string) ([]byte, error) {
	// get the table details from REST
	// setup the table API URL url
	response, err := api.Get(scopeApiURL)

	if err != nil {
		return nil, err
	}

	return response, err
}

// returns scope sys_id
func RequestScopeDataFromInstance(scopeName string) (string, error) {
	config := conf.GetConfig()

	// check if entry exists
	if exists, _, sysID := ScopeExists(scopeName); exists == true {
		err := errors.New("scope_exists")
		log.WithFields(log.Fields{"error": err}).Debug("Scope already exists, no insert!")
		return sysID, nil
	}

	// fields required here
	fields := []string{"scope", "sys_id"}

	query := fmt.Sprintf("sysparm_query=scope=%s&sysparm_fields=%s", scopeName, strings.Join(fields, ","))

	endpoint := config.GetString("app.rest.url") + "/api/now/table/sys_scope?" + query

	response, err := RequestScopeData(endpoint)

	if err != nil {
		return "", err
	}

	// unmarshal response
	var responseResult map[string]interface{}
	err = json.Unmarshal(response, &responseResult)

	if err != nil {
		fmt.Println("There was an error while unmarshalling the response!", err)
		return "", err
	}

	result, err := dyno.Get(responseResult, "result")

	if err != nil {
		fmt.Println("Error getting the result key!", err)
		return "", err
	}

	scopeSysID := ""
	// response is an array...
	for _, res := range result.([]interface{}) {
		// get scope name
		scopeName, err = dyno.GetString(res, "scope")

		if err != nil {
			log.WithFields(log.Fields{"error": err, "key": "sys_scope.scope"}).Error("Key parsing error!")
			return "", err
		}

		scopeSysID, err = dyno.GetString(res, "sys_id")

		if err != nil {
			log.WithFields(log.Fields{"error": err, "key": "sys_scope.scope"}).Error("Key parsing error!")
			return "", err
		}

		err = WriteScope(scopeSysID, scopeName)

		if err != nil {
			return "", err
		}
	}

	return scopeSysID, nil
}
