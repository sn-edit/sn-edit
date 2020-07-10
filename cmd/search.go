package cmd

import (
	"errors"
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
			conf.Err("Parsing error table flag!", log.Fields{"error": err}, true)
		}

		if len(tableName) == 0 {
			conf.Err("Please provide a valid table flag!", log.Fields{"error": errors.New("invalid_table_flag")}, true)
		}

		limit, err := cmd.Flags().GetInt64("limit")

		if err != nil {
			conf.Err("Parsing error limit flag!", log.Fields{"error": err}, true)
		}

		if limit == 0 {
			conf.Err("Please provide a valid limit flag!", log.Fields{"error": errors.New("invalid_limit_flag")}, true)
		}

		encodedQuery, err := cmd.Flags().GetString("encoded_query")

		if err != nil {
			conf.Err("Parsing error encoded_query flag!", log.Fields{"error": err}, true)
		}

		if len(encodedQuery) == 0 {
			conf.Err("Please provide a valid encoded_query flag!", log.Fields{"error": errors.New("invalid_encoded_query_flag")}, true)
		}

		fields, err := cmd.Flags().GetString("fields")

		if err != nil {
			conf.Err("Parsing error fields flag!", log.Fields{"error": err}, true)
		}

		// get table configuration from the config file
		tablesConfig := config.Get("app.tables").([]interface{})

		// get the fields for the table in question on the CLI
		tableConfigFields := conf.GetTableFieldNames(tablesConfig, tableName)

		uniqueKey, err := conf.GetUniqueKeyForTable(tablesConfig, tableName)

		if err != nil {
			conf.Err("Could not determine unique_key for the table, please validate your config!", log.Fields{"error": err}, true)
		}

		var fieldsSlice []string

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
			conf.Err("There was an error while making the request!", log.Fields{"error": err}, true)
		}

		log.WithFields(log.Fields{"result": string(response), "unique_key": uniqueKey}).Info("Search result!")
	},
}
