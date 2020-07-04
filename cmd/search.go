package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/sn-edit/sn-edit/api"
	"github.com/sn-edit/sn-edit/conf"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for results on the instance",
	Long: `This command allows you to search on the instance based on an encoded query.
This command only returns a JSON, unformatted from the instance. All the tables you would like
to search have to be present in the config file previous of using this command.`,
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

		limit, err := cmd.Flags().GetInt64("limit")

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Parsing error limit flag!")
			return
		}

		if limit == 0 {
			log.WithFields(log.Fields{"error": "no table name provided"}).Error("Please provide a valid table flag!")
			return
		}

		encodedQuery, err := cmd.Flags().GetString("encoded_query")

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Parsing error sys_id flag!")
			return
		}

		if len(encodedQuery) == 0 {
			log.WithFields(log.Fields{"error": "encoded_query not provided"}).Error("Please provide a valid encoded_query flag!")
			return
		}

		fields, err := cmd.Flags().GetString("fields")

		// get table configuration from the config file
		tablesConfig := config.Get("app.tables").([]interface{})

		// get the fields for the table in question on the CLI
		tableConfigFields := conf.GetTableFieldNames(tablesConfig, tableName)

		uniqueKey, err := conf.GetUniqueKeyForTable(tablesConfig, tableName)

		if err != nil {
			return
		}

		fieldsSlice := []string{}

		// if there are additional fields necessary
		// merge them with the tableconfig if something is found
		// otherwise simply set the config fields
		if fields != "" {
			fieldsSlice = strings.Split(fields, ",")
			for _, tableConfigField := range tableConfigFields {
				if !conf.ContainsField(fieldsSlice, tableConfigField) {
					fieldsSlice = append(fieldsSlice, tableConfigField)
				}
			}
		} else {
			fieldsSlice = tableConfigFields
		}

		recordLimit := strconv.FormatInt(limit, 10)

		searchURL := config.GetString("app.core.rest.url") + "/api/now/table/" + tableName + "?sysparm_query=" + encodedQuery + "&sysparm_fields=" + strings.Join(fieldsSlice, ",") + "&sysparm_limit=" + recordLimit

		log.WithFields(log.Fields{"url": searchURL}).Debug("Requesting url!")

		response, err := api.Get(searchURL)

		if err != nil {
			return
		}

		log.WithFields(log.Fields{"result": string(response), "unique_key": uniqueKey}).Info("Search result!")
	},
}
