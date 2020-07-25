package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sn-edit/sn-edit/api"
	"github.com/sn-edit/sn-edit/conf"
	"github.com/sn-edit/sn-edit/db"
	"github.com/sn-edit/sn-edit/file"
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

		fields, err := cmd.Flags().GetString("fields")

		if err != nil {
			conf.Err("Parsing error fields flag!", log.Fields{"error": err}, true)
		}

		fieldsSlice := strings.Split(fields, ",")

		// if array length is 1, but the first element is an empty string, do not allow processing
		if len(fieldsSlice) == 1 && fieldsSlice[0] == "" {
			conf.Err("Please provide a valid fields flag!", log.Fields{"error": errors.New("invalid_fields_flag")}, true)
		}

		// get the update set name if exists
		updateSet, err := cmd.Flags().GetString("update_set")

		if err != nil {
			conf.Err("Parsing error update_set flag!", log.Fields{"error": err}, true)
		}

		if len(updateSet) != 32 {
			log.Info("Get a list of sys_id's by calling the updateset --list command!")
			conf.Err("Please provide a valid update_set flag!", log.Fields{"error": errors.New("invalid_sys_id")}, true)
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
			conf.Err("Could not find unique_key!", log.Fields{"error": errors.New("unique_key_not_found"), "table_name": tableName, "sys_id": sysID}, true)
		}

		success, fileScopeName := db.GetEntryScopeName(tableName, sysID)

		if !success {
			conf.Err("Could not find scope for entry! Please re-download entry!", log.Fields{"error": errors.New("data_out_of_sync"), "table_name": tableName, "sys_id": sysID}, true)
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
				conf.Err("File read error! Please check permissions!", log.Fields{"error": err}, true)
			}

			data[cliField] = string(content)
		}

		// marshal into JSON
		dataJSON, err := json.Marshal(data)

		if err != nil {
			conf.Err("JSON marshalling error!", log.Fields{"error": err}, true)
		}

		// setup the upload url
		uploadURLv2 := fmt.Sprintf("%s/api/now/table/%s/%s?sysparm_fields=%s&sysparm_scope=%s", config.GetString("app.core.rest.url"), tableName, sysID, strings.Join(configFields, ","), fileScopeName)

		// if there is an update set passed
		if len(updateSet) == 32 {
			uploadURLv2 = uploadURLv2 + "&sysparm_transaction_update_set=" + updateSet
		}

		log.WithFields(log.Fields{"sys_id": sysID, "table": tableName, "fields": fieldsSlice, "scope": fileScopeName}).Info("Uploading data to the instance...")

		_, err = api.Put(uploadURLv2, dataJSON)

		if err != nil {
			conf.Err("There was an error while uploading the entry data!", log.Fields{"error": err, "sys_id": sysID, "table": tableName, "fields": fieldsSlice, "scope": fileScopeName}, true)
		}

		log.WithFields(log.Fields{"sys_id": sysID, "table": tableName, "fields": fieldsSlice, "scope": fileScopeName}).Info("The data was successfully uploaded!")
	},
}
