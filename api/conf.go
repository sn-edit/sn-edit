package api

import (
	"github.com/0x111/sn-edit/conf"
	"github.com/0x111/sn-edit/xor"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"os"
)

// load the rest credentials like username and password from the config file
// we have additional masking for password field, we need to handle this
func loadCredentials() (string, string) {
	var err error
	config := conf.GetConfig()

	username := config.GetString("app.core.rest.user")
	passwordCredential := config.GetString("app.core.rest.password")
	xorKey := config.GetString("app.core.rest.xor_key")
	isMasked := config.GetBool("app.core.rest.masked")

	if isMasked == true {
		passwordCredential = xor.EncryptDecrypt(passwordCredential, xorKey)
	}

	if isMasked == false {
		xorPassword := xor.EncryptDecrypt(config.GetString("app.core.rest.password"), config.GetString("app.core.rest.xor_key"))
		config.Set("app.core.rest.password", xorPassword)
		config.Set("app.core.rest.masked", true)
		err = config.WriteConfig()

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("There was a problem while rewriting the config file! Check the permissions please!")
			os.Exit(1)
		}
	}

	return username, passwordCredential
}

func SetupClient() {
	// Create a Resty Client
	client := resty.New()
	// load the credentials from the config file
	username, password := loadCredentials()
	// set basic auth, so every request using this client
	// will have the username and password set
	client.SetBasicAuth(username, password)
	// we shall communicate with JSON if not stated otherwise
	client.SetHeader("Content-Type", "application/json; charset=utf-8").SetHeader("Accept", "application/json")
	// set the configured client to re-use this throughout the app
	conf.SetClient(client)
}
