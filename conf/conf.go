package conf

import (
	"fmt"
	"github.com/go-resty/resty/v2"
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
			fmt.Println("Invalid config file detected!")
			fmt.Printf("The key %s is required, but could not be found! See the sample file for reference!\n", key)
			os.Exit(1)
		}
	}
}
