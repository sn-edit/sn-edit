package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/icza/dyno"
	log "github.com/sirupsen/logrus"
	"github.com/sn-edit/sn-edit/api"
	"github.com/sn-edit/sn-edit/conf"
	"strings"
)

func WriteScope(sysID string, scopeName string) error {
	//config := conf.GetConfig()
	dbc := conf.GetDB()

	stmt, err := dbc.Prepare("INSERT INTO entry_scope(sys_id, name) VALUES(?,?)")
	defer stmt.Close()

	if err != nil {
		conf.Err("There was an error while preparing the query!", log.Fields{"error": err}, false)
		return err
	}

	_, err = stmt.Exec(sysID, scopeName)

	if err != nil {
		conf.Err("There was an error while executing the query!", log.Fields{"error": err}, false)
		return err
	}

	return nil
}

func QueryScope(sysID string) (bool, int64) {
	dbc := conf.GetDB()
	stmt, err := dbc.Prepare("SELECT id FROM entry_scope WHERE sys_id=? LIMIT 1")
	defer stmt.Close()

	if err != nil {
		conf.Err("There was an error while querying the database!", log.Fields{"error": err}, false)
		return false, 0
	}

	id := int64(0)
	err = stmt.QueryRow(sysID).Scan(&id)

	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Debug("The scope was not found in the database!")
		if err == sql.ErrNoRows {
			// no rows found, it does not exist
			return false, 0
		} else {
			conf.Err("There was an error while querying the database!", log.Fields{"error": err}, false)
			return false, 0
		}
	}

	return true, id
}

func ScopeExists(scopeName string) (bool, string, string) {
	dbc := conf.GetDB()
	stmt, err := dbc.Prepare("SELECT id,sys_id FROM entry_scope WHERE name=? LIMIT 1")
	defer stmt.Close()

	if err != nil {
		conf.Err("There was an error while querying the database!", log.Fields{"error": err}, false)
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
			conf.Err("There was an error while querying the database!", log.Fields{"error": err}, false)
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
		conf.Err("There was an error while querying the database!", log.Fields{"error": err}, false)
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
			conf.Err("There was an error while querying the database!", log.Fields{"error": err}, false)
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
func RequestScopeDataFromInstance(sysScopeSysID string) (string, error) {
	config := conf.GetConfig()

	// fields required here
	fields := []string{"scope", "sys_id"}

	query := fmt.Sprintf("sysparm_query=sys_id=%s&sysparm_fields=%s", sysScopeSysID, strings.Join(fields, ","))

	endpoint := config.GetString("app.core.rest.url") + "/api/now/table/sys_scope?" + query

	log.WithFields(log.Fields{"endpoint": endpoint}).Debug("Requesting scope data")

	response, err := RequestScopeData(endpoint)

	if err != nil {
		return "", err
	}

	// unmarshal response
	var responseResult map[string]interface{}
	err = json.Unmarshal(response, &responseResult)

	if err != nil {
		conf.Err("Error unmarshalling JSON data!", log.Fields{"error": err}, false)
		return "", err
	}

	result, err := dyno.Get(responseResult, "result")

	if err != nil {
		conf.Err("Invalid key!", log.Fields{"error": err}, false)
		return "", err
	}

	scopeSysID := ""
	// response is an array...
	for _, res := range result.([]interface{}) {
		// get scope name
		scopeName, err := dyno.GetString(res, "scope")

		if err != nil {
			conf.Err("Invalid key!", log.Fields{"error": err}, false)
			return "", err
		}

		scopeSysID, err = dyno.GetString(res, "sys_id")

		if err != nil {
			conf.Err("Invalid key!", log.Fields{"error": err}, false)
			return "", err
		}

		err = WriteScope(scopeSysID, scopeName)

		if err != nil {
			return "", err
		}
	}

	return scopeSysID, nil
}
