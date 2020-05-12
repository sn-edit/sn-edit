package conf

import (
	"github.com/icza/dyno"
	log "github.com/sirupsen/logrus"
	"os"
)

func GetTableNames(tablesConfig []interface{}) []string {
	var result []string
	// iterate tables from the configuration file
	for _, value := range tablesConfig {
		tableName, err := dyno.GetString(value, "name")

		if err != nil {
			log.WithFields(log.Fields{"error": err, "path": "tablenames.name"}).Error("Was not able to find the key!")
			return []string{}
		}

		result = append(result, tableName)
	}

	return result
}

func GetTableFieldNames(tablesConfig []interface{}, tableName string) []string {
	var result []string
	// iterate tables from the configuration file
	for _, value := range tablesConfig {
		// get the table name from this and filter based on that
		table, err := dyno.GetString(value, "name")

		if err != nil {
			log.WithFields(log.Fields{"error": err, "path": "tableFieldNames.table"}).Error("Was not able to find the key!")
		}

		if table == tableName {
			// get the fields
			fields, err := dyno.Get(value, "fields")

			if err != nil {
				log.WithFields(log.Fields{"error": err, "path": "tableFieldNames.table.fields"}).Error("Was not able to find the key!")
			}

			fieldMap := fields.([]interface{})
			// now select the field names only
			for _, field := range fieldMap {
				fieldName, err := dyno.GetString(field, "field")

				if err != nil {
					log.WithFields(log.Fields{"error": err, "path": "tableFieldNames.table.field"}).Error("Was not able to find the key!")
				}

				result = append(result, fieldName)
			}
		}
	}

	return result
}

func GetFieldExtension(tablesConfig []interface{}, tableName string, fieldName string) string {
	// iterate tables from the configuration file
	for _, value := range tablesConfig {
		// get the table name from this and filter based on that
		table, err := dyno.GetString(value, "name")

		if err != nil {
			log.WithFields(log.Fields{"error": err, "path": "fieldExtension.name"}).Error("Was not able to find the key!")
			os.Exit(1)
		}

		if table == tableName {
			// get the fields
			fields, err := dyno.Get(value, "fields")

			if err != nil {
				log.WithFields(log.Fields{"error": err, "path": "fieldExtension.fields"}).Error("Was not able to find the key!")
				os.Exit(1)
			}

			fieldMap := fields.([]interface{})
			// now select the field names only
			for _, field := range fieldMap {
				fieldNeedle, err := dyno.GetString(field, "field")

				if err != nil {
					log.WithFields(log.Fields{"error": err, "path": "fieldExtension.fields.field"}).Error("Was not able to find the key!")
					os.Exit(1)
				}

				if fieldNeedle == fieldName {
					fieldExtension, err := dyno.GetString(field, "extension")

					if err != nil {
						log.WithFields(log.Fields{"error": err, "path": "fieldExtension.fields.field"}).Error("Was not able to find the key!")
						os.Exit(1)
					}

					return fieldExtension
				}
			}
		}
	}

	return ""
}

func GetUniqueKeyForTable(tablesConfig []interface{}, tableName string) (string, error) {
	var uniqueKey string
	// iterate tables from the configuration file
	for _, value := range tablesConfig {
		// get the table name from this and filter based on that
		table, err := dyno.GetString(value, "name")

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Was not able to find the key!")
			os.Exit(1)
		}

		if table == tableName {
			// get the fields
			uniqueKey, err = dyno.GetString(value, "unique_key")

			if err != nil {
				log.WithFields(log.Fields{"error": err}).Error("Was not able to find the key!")
				os.Exit(1)
			}
		}
	}
	return uniqueKey, nil
}

func EnforceFields(tablesConfig []interface{}, tableName string, fields []string) []string {
	uniqueKey, err := GetUniqueKeyForTable(tablesConfig, tableName)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Please define a unique key for every table for sn-edit!")
		os.Exit(1)
	}

	requiredFields := []string{"sys_id", "sys_scope.name", uniqueKey}
	for _, requiredField := range requiredFields {
		if !containsField(fields, requiredField) {
			fields = append(fields, requiredField)
		}
	}

	return fields
}

func containsField(fields []string, key string) bool {
	for _, field := range fields {
		if field == key {
			return true
		}
	}

	return false
}
