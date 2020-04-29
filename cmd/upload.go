package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/0x111/sn-edit/api"
	"github.com/0x111/sn-edit/conf"
	"github.com/0x111/sn-edit/directory"
	"github.com/0x111/sn-edit/file"
	"github.com/icza/dyno"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var uploadEntryCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload one entry from servicenow",
	Long: `You can upload one entry (for example a script) to the instance.
Provide a table name, sys_id and field please. The table name and fields should be already configured in the config file.
Otherwise sn-edit will not be able to determine the location or download the data to.
Providing a field is optional, if you do not provide any, sn-edit will assume you would like to update the contents of every field for the entry saved locally.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := conf.GetConfig()

		tableName, err := cmd.Flags().GetString("table")

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Parsing error table name flag!")
			return
		}

		if len(tableName) == 0 {
			log.WithFields(log.Fields{"error": "no table name provided"}).Error("Please provide a valid table flag!")
			return
		}

		sysID, err := cmd.Flags().GetString("sys_id")

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Parsing error sys_id flag!")
			return
		}

		if len(sysID) != 32 {
			log.WithFields(log.Fields{"error": "sys_id length must be 32"}).Error("Please provide a valid sys_id flag!")
			return
		}

		fields, err := cmd.Flags().GetString("fields")

		log.WithFields(log.Fields{"sys_id": sysID, "table": tableName}).Info("Uploading data to the instance")

		fieldsSlice := strings.Split(fields, ",")

		// get table configuration from the config file
		tablesConfig := config.Get("app.tables").([]interface{})

		// get the fields for the table in question on the CLI
		//fields := conf.GetTableFieldNames(tablesConfig, tableName)

		// setup the download url
		downloadURL := config.GetString("app.rest.url") + "/api/now/table/" + tableName + "/" + sysID + "?sysparm_fields=" + strings.Join(fieldsSlice, ",")

		response, err := api.Get(downloadURL)

		if err != nil {
			return
		}

		// unmarshal response
		var responseResult map[string]interface{}
		err = json.Unmarshal(response, &responseResult)

		if err != nil {
			fmt.Println("There was an error while unmarshaling the response!", err)
			return
		}

		result, err := dyno.Get(responseResult, "result")

		if err != nil {
			fmt.Println("Error getting the result key!", err)
			return
		}

		// iterate through the entries
		//for _, entry := range results {
		fieldSysName, err := dyno.GetString(result, "sys_name")

		if err != nil {
			log.WithFields(log.Fields{"error": err, "key": "sys_name"}).Error("There was an error while getting the unique key!")
		}

		// create directory for sys_name
		directoryPath := config.GetString("app.root_directory") + string(os.PathSeparator) + tableName + string(os.PathSeparator) + fieldSysName
		_, err = directory.CreateDirectoryStructure(directoryPath)

		if err != nil {
			log.WithFields(log.Fields{"error": err, "directory": directoryPath}).Error("There was an error while creating the directory structure! Please check the permissions!")
			return
		}

		// todo error handling
		// go through all the fields that are defined in the config
		for _, fieldName := range fields {
			fieldContent, err := dyno.GetString(result, fieldName)

			if err != nil {
				log.WithFields(log.Fields{"error": err, "key": fieldName}).Error("There was an error while getting the key!")
			}

			fieldExtension := conf.GetFieldExtension(tablesConfig, tableName, fieldName)

			err = file.WriteFile(tableName, fieldSysName, fieldName, fieldExtension, []byte(fieldContent))

			if err != nil {
				log.WithFields(log.Fields{"error": err}).Error("File writing error! Check permissions please!")
			}
		}
	},
}
