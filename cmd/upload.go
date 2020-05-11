package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/0x111/sn-edit/api"
	"github.com/0x111/sn-edit/conf"
	"github.com/0x111/sn-edit/db"
	"github.com/0x111/sn-edit/file"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

		fieldsSlice := strings.Split(fields, ",")

		// if array length is 1, but the first element is an empty string, do not allow processing
		if len(fieldsSlice) == 1 && fieldsSlice[0] == "" {
			log.WithFields(log.Fields{"error": "please provide the fields"}).Error("Please provide a valid fields flag!")
			return
		}

		// get the update set name if exists
		updateSet, err := cmd.Flags().GetString("update_set")

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Parsing error update_set flag!")
			return
		}

		if len(sysID) != 32 {
			log.WithFields(log.Fields{"error": "update_set length must be 32"}).Error("Please provide a valid update_set flag!")
			return
		}

		// get table configuration from the config file
		tablesConfig := config.Get("app.tables").([]interface{})

		// get the fields for the table in question on the CLI
		configFields := conf.GetTableFieldNames(tablesConfig, tableName)

		// todo: check if valid fields provided, compare tableconfig with the fieldsSlice

		// build data
		data := make(map[string]interface{})

		found, uniqueKeyName := db.QueryUniqueKey(tableName, sysID)

		if !found {
			log.WithFields(log.Fields{"error": err, "table_name": tableName, "sys_id": sysID}).Error("There was an error while getting the key!")
			return
		}

		success, fileScopeName := db.GetEntryScopeName(tableName, sysID)

		if !success {
			log.WithFields(log.Fields{"error": err, "table_name": tableName, "sys_id": sysID}).Error("Please re-download the entry again!")
			return
		}

		// iterate through the cli fields which need updating on the instance
		for _, cliField := range fieldsSlice {
			// find the extension based on the field and tableName from the config
			extension := conf.GetFieldExtension(tablesConfig, tableName, cliField)
			// generate file path to the given file, full path, to read
			filePath := file.GenerateFilePath(tableName, fileScopeName, uniqueKeyName, cliField, extension)
			// get the contents of the file
			content, err := file.ReadFile(filePath)

			if err != nil {
				log.WithFields(log.Fields{"error": err, "file_path": filePath}).Error("There was an error while getting the key!")
				return
			}

			data[cliField] = string(content)
		}

		// marshal into JSON
		dataJSON, err := json.Marshal(data)

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Error while marshalling JSON data!")
			return
		}

		// setup the upload url
		uploadURLv2 := fmt.Sprintf("%s/api/now/table/%s/%s?sysparm_fields=%s&sysparm_scope=%s", config.GetString("app.core.rest.url"), tableName, sysID, strings.Join(configFields, ","), fileScopeName)

		// if there is an update set passed
		if len(updateSet) == 32 {
			uploadURLv2 = uploadURLv2 + "&sysparm_transaction_update_set=" + updateSet
		}

		log.WithFields(log.Fields{"sys_id": sysID, "table": tableName, "fields": fieldsSlice, "scope": fileScopeName}).Info("Uploading data to the instance")

		_, err = api.Put(uploadURLv2, dataJSON)

		if err != nil {
			log.WithFields(log.Fields{"error": err, "sys_id": sysID, "table": tableName, "fields": fieldsSlice, "scope": fileScopeName}).Error("There was an error while uploading the data!")
			return
		}

		log.WithFields(log.Fields{"sys_id": sysID, "table": tableName, "fields": fieldsSlice, "scope": fileScopeName}).Info("The data was successfully uploaded!")
	},
}
