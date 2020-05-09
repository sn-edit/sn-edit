package updateset

import (
	"encoding/json"
	"fmt"
	"github.com/0x111/sn-edit/api"
	"github.com/0x111/sn-edit/conf"
	"github.com/0x111/sn-edit/db"
	log "github.com/sirupsen/logrus"
)

func SetCommand(scopeName string, updateSetSysID string) {
	config := conf.GetConfig()

	found, scopeID, _ := db.ScopeExists(scopeName)

	if !found {
		log.WithFields(log.Fields{"error": "scope_not_found", "found": false, "scope_name": scopeName}).Error("Could not find scope in the DB!")
		return
	}

	// build data
	data := make(map[string]interface{})

	found, sysID, name := db.QueryUpdateSet(updateSetSysID)

	if !found {
		log.WithFields(log.Fields{"error": "update_set_not_found", "found": false, "scope_name": scopeName, "update_set": updateSetSysID}).Error("Could not find Update set in the DB!")
		return
	}

	data["sysId"] = sysID
	data["name"] = name

	dataJSON, err := json.Marshal(data)

	if err != nil {
		fmt.Println("Error marshalling json!", err)
		return
	}

	log.Infof("Setting your Update Set in the scope %s to %s!", scopeName, name)

	// make request to the instance (to get an updated list of scopes for the scope in the CLI)
	setUpdateSetEndPoint := config.GetString("app.core.rest.url") + "/api/now/ui/concoursepicker/updateset?sysparm_transaction_scope=" + scopeID
	_, err = api.Put(setUpdateSetEndPoint, dataJSON)

	if err != nil {
		fmt.Println("ERROR", err)
		return
	}

	log.Infof("The Update Set was updated in the %s scope to %s!", scopeName, name)
}
