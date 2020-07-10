package conf

import (
	"github.com/icza/dyno"
	log "github.com/sirupsen/logrus"
)

func GetTableNames(tablesConfig []interface{}) []string {
	var result []string
	// iterate tables from the configuration file
	for _, value := range tablesConfig {
		tableName, err := dyno.GetString(value, "name")

		if err != nil {
			Err("Invalid key!", log.Fields{"error": err}, false)
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
			Err("Invalid key!", log.Fields{"error": err}, false)
		}

		if table == tableName {
			// get the fields
			fields, err := dyno.Get(value, "fields")

			if err != nil {
				Err("Invalid key!", log.Fields{"error": err}, false)
			}

			fieldMap := fields.([]interface{})
			// now select the field names only
			for _, field := range fieldMap {
				fieldName, err := dyno.GetString(field, "field")

				if err != nil {
					Err("Invalid key!", log.Fields{"error": err}, false)
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
			Err("Invalid key!", log.Fields{"error": err}, true)
		}

		if table == tableName {
			// get the fields
			fields, err := dyno.Get(value, "fields")

			if err != nil {
				Err("Invalid key!", log.Fields{"error": err}, true)
			}

			fieldMap := fields.([]interface{})
			// now select the field names only
			for _, field := range fieldMap {
				fieldNeedle, err := dyno.GetString(field, "field")

				if err != nil {
					Err("Invalid key!", log.Fields{"error": err}, true)
				}

				if fieldNeedle == fieldName {
					fieldExtension, err := dyno.GetString(field, "extension")

					if err != nil {
						Err("Invalid key!", log.Fields{"error": err}, true)
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
			Err("Invalid key!", log.Fields{"error": err}, true)
		}

		if table == tableName {
			// get the fields
			uniqueKey, err = dyno.GetString(value, "unique_key")

			if err != nil {
				Err("Invalid key!", log.Fields{"error": err}, true)
			}
		}
	}
	return uniqueKey, nil
}

func EnforceFields(tablesConfig []interface{}, tableName string, fields []string) []string {
	uniqueKey, err := GetUniqueKeyForTable(tablesConfig, tableName)

	if err != nil {
		Err("Please define a unique key for every table in your cocnfig!", log.Fields{"error": err}, true)
	}

	requiredFields := []string{"sys_id", "sys_scope.name", "sys_scope.sys_id", uniqueKey}
	for _, requiredField := range requiredFields {
		if !ContainsField(fields, requiredField) {
			fields = append(fields, requiredField)
		}
	}

	return fields
}

func ContainsField(fields []string, key string) bool {
	for _, field := range fields {
		if field == key {
			return true
		}
	}

	return false
}
