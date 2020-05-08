package updateset

import (
	"encoding/json"
	"fmt"
	"github.com/0x111/sn-edit/api"
	"github.com/0x111/sn-edit/conf"
	"github.com/0x111/sn-edit/db"
	"github.com/icza/dyno"
	log "github.com/sirupsen/logrus"
)

func ListCommand(scopeName string) {
	config := conf.GetConfig()

	scopeSysID, err := db.RequestScopeDataFromInstance(scopeName)

	if err != nil {
		return
	}

	found, scopeID := db.QueryScope(scopeSysID)

	if !found {
		log.WithFields(log.Fields{"error": "scope_not_found", "found": false, "scope_name": scopeName}).Error("Could not find scope in the DB!")
		return
	}

	// make request to the instance (to get an updated list of scopes for the scope in the CLI)
	listUpdateSetEndpoint := config.GetString("app.rest.url") + "/api/now/ui/concoursepicker/updateset?sysparm_transaction_scope=" + scopeSysID
	response, err := api.Get(listUpdateSetEndpoint)

	if err != nil {
		fmt.Println("ERROR", err)
		return
	}

	// unmarshal response
	var responseResult map[string]interface{}
	err = json.Unmarshal(response, &responseResult)

	if err != nil {
		fmt.Println("There was an error while unmarshalling the response!", err)
		return
	}

	result, err := dyno.Get(responseResult, "result", "updateSet")

	if err != nil {
		fmt.Println("Error getting the result key!", err)
		return
	}

	updateSets := result.([]interface{})

	// print current update set
	currentUpdateSet, err := dyno.Get(responseResult, "result", "current")

	if err != nil {
		log.WithFields(log.Fields{"error": err, "path": "updateset.result.current"}).Error("Was not able to parse the response correctly!")
		return
	}

	currentName, err := dyno.GetString(currentUpdateSet, "name")

	if err != nil {
		log.WithFields(log.Fields{"error": err, "path": "updateset.name"}).Error("Was not able to parse the response correctly!")
		return
	}

	currentSysID, err := dyno.GetString(currentUpdateSet, "sysId")

	if err != nil {
		log.WithFields(log.Fields{"error": err, "path": "updateset.sysId"}).Error("Was not able to parse the response correctly!")
		return
	}

	fmt.Printf("Currently selected update set for the %s scope\n", scopeName)

	fmt.Printf("Update Set: %s\n", currentName)
	fmt.Printf("Sys id: %s\n", currentSysID)
	fmt.Println("------------------------------")
	fmt.Printf("List of Update sets for %s scope\n", scopeName)

	for _, updateSet := range updateSets {
		sysID, err := dyno.GetString(updateSet, "sysId")

		if err != nil {
			log.WithFields(log.Fields{"error": err, "key": "updateSet.sysId"}).Error("There was an error while getting the key!")
		}

		name, err := dyno.GetString(updateSet, "name")

		if err != nil {
			log.WithFields(log.Fields{"error": err, "key": "updateSet.name"}).Error("There was an error while getting the key!")
		}

		// write to db if not exists
		err = db.WriteUpdateSet(name, sysID, scopeID)

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Debug("There was an error while writing the Update Set data!")
			return
		}

		if sysID == currentSysID {
			continue
		}

		fmt.Printf("Update set: %s\n", name)
		fmt.Printf("Sys id: %s\n", sysID)
	}
}
