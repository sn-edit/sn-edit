package conf

import (
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
