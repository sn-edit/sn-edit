package conf

import (
	"github.com/icza/dyno"
)

// todo: error handling
func GetTableNames(tablesConfig []interface{}) []string {
	var result []string
	// iterate tables from the configuration file
	for _, value := range tablesConfig {
		tableName, _ := dyno.GetString(value, "name")
		result = append(result, tableName)
	}

	return result
}

// todo: error handling
func GetTableFieldNames(tablesConfig []interface{}, tableName string) []string {
	var result []string
	// iterate tables from the configuration file
	for _, value := range tablesConfig {
		// get the table name from this and filter based on that
		table, _ := dyno.GetString(value, "name")
		if table == tableName {
			// get the fields
			fields, _ := dyno.Get(value, "fields")
			fieldMap := fields.([]interface{})
			// now select the field names only
			for _, field := range fieldMap {
				fieldName, _ := dyno.GetString(field, "field")
				result = append(result, fieldName)
			}
		}
	}

	return result
}

// todo: error handling
func GetFieldExtension(tablesConfig []interface{}, tableName string, fieldName string) string {
	// iterate tables from the configuration file
	for _, value := range tablesConfig {
		// get the table name from this and filter based on that
		table, _ := dyno.GetString(value, "name")
		if table == tableName {
			// get the fields
			fields, _ := dyno.Get(value, "fields")
			fieldMap := fields.([]interface{})
			// now select the field names only
			for _, field := range fieldMap {
				fieldNeedle, _ := dyno.GetString(field, "field")
				if fieldNeedle == fieldName {
					fieldExtension, _ := dyno.GetString(field, "extension")
					return fieldExtension
				}
			}
		}
	}

	return ""
}
