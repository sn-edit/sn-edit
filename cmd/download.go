package cmd

import (
	"encoding/json"
	"errors"
	"github.com/icza/dyno"
	log "github.com/sirupsen/logrus"
	"github.com/sn-edit/sn-edit/api"
	"github.com/sn-edit/sn-edit/conf"
	"github.com/sn-edit/sn-edit/db"
	"github.com/sn-edit/sn-edit/directory"
	"github.com/sn-edit/sn-edit/file"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var downloadEntryCmd = &cobra.Command{
	Use:   "download",
	Short: "Download one entry from servicenow",
	Long: `You can download one entry (for example a script) from the instance.
Provide a table name and sys_id please. The table name and fields should be already configured in the config file.
Otherwise sn-edit will not be able to determine the location or download the data to.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := conf.GetConfig()

		tableName, err := cmd.Flags().GetString("table")

		if err != nil {
			conf.Err("Parsing error table flag!", log.Fields{"error": err}, true)
		}

		if len(tableName) == 0 {
			conf.Err("Please provide a valid table flag!", log.Fields{"error": errors.New("invalid_table")}, true)
		}

		sysID, err := cmd.Flags().GetString("sys_id")

		if err != nil {
			conf.Err("Parsing error sys_id flag!", log.Fields{"error": err}, true)
		}

		if len(sysID) != 32 {
			conf.Err("Please provide a valid sys_id flag!", log.Fields{"error": errors.New("invalid_sys_id")}, true)
		}

		// get table configuration from the config file
		tablesConfig := config.Get("app.tables").([]interface{})

		// get the fields for the table in question on the CLI
		fields := conf.GetTableFieldNames(tablesConfig, tableName)

		// enforce sys_id and scope if not present already
		fields = conf.EnforceFields(tablesConfig, tableName, fields)

		// setup the download url
		downloadURL := config.GetString("app.core.rest.url") + "/api/now/table/" + tableName + "/" + sysID + "?sysparm_fields=" + strings.Join(fields, ",")

		log.WithFields(log.Fields{"api_url": downloadURL}).Debug()
		log.WithFields(log.Fields{"sys_id": sysID, "table": tableName, "fields": fields}).Info("Downloading the data from the instance")

		response, err := api.Get(downloadURL)

		if err != nil {
			return
		}

		// unmarshal response
		var responseResult map[string]interface{}
		err = json.Unmarshal(response, &responseResult)

		if err != nil {
			conf.Err("There was an error while unmarshalling the response!", log.Fields{"error": err}, true)
		}

		result, err := dyno.Get(responseResult, "result")

		if err != nil {
			conf.Err(err, log.Fields{"error": err}, true)
		}

		uniqueKey, err := conf.GetUniqueKeyForTable(tablesConfig, tableName)

		if err != nil {
			conf.Err("Invalid tables config!", log.Fields{"error": err}, true)
		}

		uniqueKeyName, err := dyno.GetString(result, uniqueKey)

		if err != nil {
			conf.Err("Invalid unique key!", log.Fields{"error": err}, true)
		}

		log.WithFields(log.Fields{"name": uniqueKeyName}).Debug("Entry identified!")

		fieldScopeSysID, err := dyno.GetString(result, "sys_scope.sys_id")

		if err != nil {
			// if no scope found, fallback to global
			fieldScopeSysID = "global"
			conf.Err("Invalid scope for entry!", log.Fields{"error": err}, true)
		}

		// write entry to the db
		err = db.WriteEntry(tableName, uniqueKeyName, sysID, fieldScopeSysID)

		if err != nil {
			conf.Err("Could not write entry to the database!", log.Fields{"error": err}, true)
		}

		found, fieldScopeName := db.GetScopeNameFromSysID(fieldScopeSysID)

		// scope names in lowercase in folder structure
		fieldScopeName = strings.ToLower(fieldScopeName)

		if !found {
			conf.Err("Scope not found in the database!", log.Fields{"error": err, "name": fieldScopeName, "sys_id": fieldScopeSysID}, true)
		}

		// create directory for sys_name
		directoryPath := config.GetString("app.core.root_directory") + string(os.PathSeparator) + fieldScopeName + string(os.PathSeparator) + tableName + string(os.PathSeparator) + file.FilterSpecialChars(uniqueKeyName)
		_, err = directory.CreateDirectoryStructure(directoryPath)

		if err != nil {
			conf.Err("Error while creating directory structure!", log.Fields{"error": err, "directory": directoryPath}, true)
		}

		// go through all the fields that are defined in the config
		for _, fieldName := range fields {
			// we do not need to download sys_scope
			if strings.Contains(fieldName, "scope") {
				continue
			}

			fieldContent, err := dyno.GetString(result, fieldName)

			if err != nil {
				conf.Err("Invalid key!", log.Fields{"error": err}, true)
			}

			fieldExtension := conf.GetFieldExtension(tablesConfig, tableName, fieldName)

			err = file.WriteFile(tableName, fieldScopeName, uniqueKeyName, fieldName, fieldExtension, []byte(fieldContent))

			if err != nil {
				conf.Err("File write error! Please check permissions!", log.Fields{"error": err}, true)
			}
		}

		log.WithFields(log.Fields{"table_name": tableName, "sys_id": sysID}).Info("Entry successfully downloaded!")
	},
}
