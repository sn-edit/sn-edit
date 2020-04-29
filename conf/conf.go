package conf

import(
	"github.com/0x111/sn-edit/directory"
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
	switch conf.GetString("app.log_level") {
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

func SetupDirectoryStructure() {
	config := GetConfig()

	// get table configuration from the config file
	tablesConfig := config.Get("app.tables").([]interface{})
	// get all the table names
	tablesData := GetTableNames(tablesConfig)

	for _, tableName := range tablesData {
		// create directory structure, main table folders
		_, err := directory.CreateDirectoryStructure(config.GetString("app.root_directory") + string(os.PathSeparator) + tableName)

		if err != nil {
			log.WithFields(log.Fields{"path": config.GetString("app.root_directory") + string(os.PathSeparator) + tableName}).Error("We could not create the directory under the provided path!")
		}
	}

}