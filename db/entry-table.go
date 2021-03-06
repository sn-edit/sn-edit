package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/icza/dyno"
	log "github.com/sirupsen/logrus"
	"github.com/sn-edit/sn-edit/api"
	"github.com/sn-edit/sn-edit/conf"
)

func WriteTable(tableName string) error {
	config := conf.GetConfig()
	dbc := conf.GetDB()

	// check if table exists
	if exists, _ := TableExists(tableName); exists == true {
		log.WithFields(log.Fields{"error": "table_exists"}).Debug("Table already exists, no insert!")
		return nil
	}

	// get the table details from REST
	// setup the table API URL url
	// https://devxxxx.service-now.com/api/now/table/sys_db_object?sysparm_query=name=sys_db_object&sysparm_fields=sys_id,sys_scope,name&sysparm_limit=1
	tableAPIURL := config.GetString("app.core.rest.url") + "/api/now/table/sys_db_object?sysparm_query=name=" + tableName + "&sysparm_fields=sys_id,sys_scope.sys_id,sys_scope.name,name&sysparm_limit=1"

	response, err := api.Get(tableAPIURL)

	if err != nil {
		return err
	}

	// unmarshal response
	var responseResult map[string]interface{}
	err = json.Unmarshal(response, &responseResult)

	if err != nil {
		conf.Err("Error unmarshalling JSON data!", log.Fields{"error": err}, false)
		return err
	}

	result, err := dyno.Get(responseResult, "result")

	if err != nil {
		conf.Err("Invalid key!", log.Fields{"error": err}, false)
		return err
	}

	// scopeDataID needed
	scopeDataID := ""
	resultTableName := ""
	resultTableSysID := ""

	for _, res := range result.([]interface{}) {
		// table name
		resultTableName, err = dyno.GetString(res, "name")

		if err != nil {
			conf.Err("Invalid key!", log.Fields{"error": err}, false)
			return err
		}

		// table sys_id
		resultTableSysID, err = dyno.GetString(res, "sys_id")

		if err != nil {
			conf.Err("Invalid key!", log.Fields{"error": err}, false)
			return err
		}

		// table scope
		scopeName, err := dyno.GetString(res, "sys_scope.name") // scope name

		if err != nil {
			conf.Err("Invalid key!", log.Fields{"error": err}, false)
			return err
		}

		scopeDataID, err = dyno.GetString(res, "sys_scope.sys_id") // sys_id of scope

		if err != nil {
			conf.Err("Invalid key!", log.Fields{"error": err}, false)
			return err
		}

		// check for scope and insert if non-existent
		// GET scope data for table
		err = WriteScope(scopeDataID, scopeName)

		if err != nil {
			conf.Err("Could not write scope data!", log.Fields{"error": err}, false)
			return err
		}
	}

	success, scopeID := QueryScope(scopeDataID)

	if !success {
		conf.Err("Error while querying database data!", log.Fields{"error": errors.New("scope_data_query")}, false)
		return err
	}

	stmt, err := dbc.Prepare("INSERT INTO entry_table (sys_id, name, sys_scope) VALUES(?,?,?)")
	defer stmt.Close()

	if err != nil {
		conf.Err("Error while preparing the query!", log.Fields{"error": err}, false)
		return err
	}

	_, err = stmt.Exec(resultTableSysID, resultTableName, scopeID)

	if err != nil {
		conf.Err("Error while executing the query!", log.Fields{"error": err}, false)
		return err
	}

	return nil
}

func QueryTable(tableName string) (bool, string) {
	dbc := conf.GetDB()
	stmt, err := dbc.Prepare("SELECT id FROM entry_table WHERE name=? LIMIT 1")
	defer stmt.Close()

	if err != nil {
		conf.Err("Error while querying database data!", log.Fields{"error": err}, false)
		return false, ""
	}

	id := ""
	err = stmt.QueryRow(tableName).Scan(&id)

	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Debug("The table entry was not found in the database!")
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

func TableExists(tableName string) (bool, string) {
	dbc := conf.GetDB()
	stmt, err := dbc.Prepare("SELECT id FROM entry_table WHERE name=? LIMIT 1")
	defer stmt.Close()

	if err != nil {
		conf.Err("Error while querying database data!", log.Fields{"error": err}, false)
		return false, ""
	}

	sysID := ""
	err = stmt.QueryRow(tableName).Scan(&sysID)

	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Debug("The table was not found in the database!")
		if err == sql.ErrNoRows {
			// no rows found, it does not exist
			return false, ""
		} else {
			conf.Err("Error while querying database data!", log.Fields{"error": err}, false)
			return false, ""
		}
	}

	return true, sysID
}
