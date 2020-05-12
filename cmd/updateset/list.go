package updateset

import (
	"encoding/json"
	"fmt"
	"github.com/icza/dyno"
	log "github.com/sirupsen/logrus"
	"github.com/sn-edit/sn-edit/api"
	"github.com/sn-edit/sn-edit/conf"
	"github.com/sn-edit/sn-edit/db"
	"github.com/spf13/cobra"
)

func ListCommand(cmd *cobra.Command, scopeName string) {
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
	listUpdateSetEndpoint := config.GetString("app.core.rest.url") + "/api/now/ui/concoursepicker/updateset?sysparm_transaction_scope=" + scopeSysID
	response, err := api.Get(listUpdateSetEndpoint)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Error during the request to the instance!")
		return
	}

	// unmarshal response
	var responseResult map[string]interface{}
	err = json.Unmarshal(response, &responseResult)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while unmarshalling the response!")
		return
	}

	result, err := dyno.Get(responseResult, "result", "updateSet")

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Error while finding the result key!")
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

	// build data json.. eventually
	var data = map[string]interface{}{}
	// add the current update set
	data["current"] = map[string]string{"name": currentName, "sys_id": currentSysID}
	data["others"] = []interface{}{}

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

		data["others"] = append(data["others"].([]interface{}), map[string]string{"name": name, "sys_id": sysID})
	}

	if outputJSON, _ := cmd.Flags().GetBool("json"); outputJSON {
		log.WithFields(data).Info()
	} else {
		fmt.Printf("Currently selected update set for the %s scope\n", scopeName)

		fmt.Printf("Update Set: %s\n", currentName)
		fmt.Printf("Sys id: %s\n", currentSysID)
		fmt.Println("------------------------------")
		fmt.Printf("List of Update sets for %s scope\n", scopeName)

		for _, row := range data["others"].([]interface{}) {
			name := row.(map[string]string)["name"]
			sysID := row.(map[string]string)["sys_id"]

			fmt.Print("\n")
			fmt.Printf("Update set: %s\n", name)
			fmt.Printf("Sys id: %s\n", sysID)
		}
	}
}
