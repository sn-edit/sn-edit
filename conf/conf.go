package conf

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/icza/dyno"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

var conf *viper.Viper
var restClient *resty.Client

func SetConfig(config *viper.Viper) {
	conf = config
}

func GetConfig() *viper.Viper {
	return conf
}

func SetClient(client *resty.Client) {
	restClient = client
}

func GetClient() *resty.Client {
	return restClient
}

func SetLoggerLevel() {
	switch conf.GetString("app.core.log_level") {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	}
}

// this function should validate the configuration
// checks if everything is in order
func ValidateConfig() {
	config := GetConfig()

	// list of required keys
	// if these do not exist, prevent running the app
	requiredKeys := []string{
		"app.core.root_directory",
		"app.core.log_level",
		"app.core.rest.url",
		"app.core.rest.user",
		"app.core.rest.xor_key",
		"app.core.rest.password",
		"app.core.rest.masked",
		"app.core.db.initialised",
		"app.core.db.path",
	}

	for _, key := range requiredKeys {
		if !config.IsSet(key) {
			log.WithFields(log.Fields{"error": fmt.Sprintf("The key %s is required, but could not be found! See the sample file for reference!", key)}).Error("Invalid config file detected!")
			os.Exit(1)
		}
	}
}

// every table should have some fields defined, these fields should have an extension set
func ValidateTableData() {
	config := GetConfig()

	// get table configuration from the config file
	tablesConfig := config.Get("app.tables").([]interface{})

	for _, table := range tablesConfig {
		tableName, err := dyno.GetString(table, "name")

		if err != nil {
			log.WithFields(log.Fields{"error": err, "path": "validateTableData.table_name"}).Error("Was not able to find the key!")
			os.Exit(1)
		}

		tableFields, err := dyno.Get(table, "fields")
		fields := tableFields.([]interface{})

		if err != nil {
			log.WithFields(log.Fields{"error": err, "table": tableName, "path": "validateTableData.fields"}).Error("Was not able to find the key!")
			os.Exit(1)
		}

		for _, field := range fields {
			fieldName, err := dyno.GetString(field, "field")

			if err != nil {
				log.WithFields(log.Fields{"error": err, "table": tableName, "fieldName": fieldName}).Error("Was not able to find the key!")
				os.Exit(1)
			}

			fieldExtension, err := dyno.GetString(field, "extension")

			if err != nil {
				log.WithFields(log.Fields{"error": err, "table": tableName, "fieldName": fieldName, "extension": fieldExtension}).Error("Was not able to find the key! Every field does need to have an extension!")
				os.Exit(1)
			}

			if len(fieldName) == 0 {
				log.WithFields(log.Fields{"error": "field", "table": tableName, "fieldName": fieldName, "extension": fieldExtension}).Error("Missing correct field name!")
				os.Exit(1)
			}

			if len(fieldExtension) == 0 {
				log.WithFields(log.Fields{"error": "extension", "table": tableName, "fieldName": fieldName, "extension": fieldExtension}).Error("Every field does need to have an extension!")
				os.Exit(1)
			}
		}

	}
}
