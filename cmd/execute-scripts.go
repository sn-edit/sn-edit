package cmd

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sn-edit/sn-edit/conf"
	"github.com/sn-edit/sn-edit/db"
	"github.com/sn-edit/sn-edit/file"
	"github.com/sn-edit/sn-edit/xor"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
)

var executeScriptsCmd = &cobra.Command{
	Use:   "execute",
	Short: "Run scripts on your instance",
	Long: `The command makes it easy and simple for you to run scripts on your instance
like you would with the Scripts - Background functionality. The response is returned as-is.
No formatting applied since there is no strict format anyways. You should pass a file to this
command and it would then run the contents on the instance. For this to work, you need a user
that has access to the Background Scripts functionality. Otherwise the feature may not work correctly!`,
	Run: func(cmd *cobra.Command, args []string) {
		config := conf.GetConfig()
		scriptFile, err := cmd.Flags().GetString("file")

		if err != nil {
			conf.Err("Parsing error file flag!", log.Fields{"error": err}, true)
		}

		if len(scriptFile) == 0 {
			conf.Err("Provide a valid file flag please!", log.Fields{"error": errors.New("invalid_file_flag")}, true)
		}

		scopeName, err := cmd.Flags().GetString("scope")

		if err != nil {
			conf.Err("Parsing error scope flag!", log.Fields{"error": err}, true)
		}

		scopeSysID, err := db.RequestScopeDataFromInstance(scopeName)

		if err != nil {
			conf.Err("Error while requesting scope data!", log.Fields{"error": err}, true)
		}

		if exists := file.Exists(scriptFile); exists == false {
			conf.Err("The file could not be found! Check the file path!", log.Fields{"error": err, "file": scriptFile}, true)
		}

		// read file contents
		dat, err := ioutil.ReadFile(scriptFile)

		//fmt.Println(string(dat))
		cookieJar, err := cookiejar.New(nil)

		if err != nil {
			conf.Err("Error initializing cookie jar!", log.Fields{"error": err}, true)
		}

		client := &http.Client{Jar: cookieJar}

		username := config.GetString("app.core.rest.user")
		passwordCredential := config.GetString("app.core.rest.password")
		xorKey := config.GetString("app.core.rest.xor_key")
		isMasked := config.GetBool("app.core.rest.masked")

		if isMasked == true {
			passwordCredential = xor.EncryptDecrypt(passwordCredential, xorKey)
		}

		// login to the instance to get CSRF token
		loginUrl, err := url.Parse(config.GetString("app.core.rest.url") + "/login.do")

		if err != nil {
			conf.Err("Could not parse provided URL!", log.Fields{"error": err}, true)
		}

		// get the CK key for login
		resp1, err := client.Get(loginUrl.String())

		if err != nil {
			conf.Err("There was an error while making the request!", log.Fields{"error": err}, true)
		}

		defer resp1.Body.Close()

		body, err := ioutil.ReadAll(resp1.Body)

		// search for the CK key in the HTML string
		re := regexp.MustCompile(`<input name="sysparm_ck" id="sysparm_ck" type="hidden" value="([a-z0-9]{72})`)
		match := re.FindStringSubmatch(string(body))

		if match == nil {
			conf.Err("The ck key could not be found in the provided HTML!", log.Fields{"error": errors.New("ck_key_not_found")}, true)
		}

		ckToken := match[1]

		// build the form values manually
		form := url.Values{}
		form.Add("user_name", username)
		form.Add("user_password", passwordCredential)
		form.Add("sys_action", "sysverb_login")
		form.Add("sysparm_ck", ckToken)

		// login to the instance
		resp2, err := client.PostForm(loginUrl.String(), form)

		if err != nil {
			conf.Err("There was an error with the login process!", log.Fields{"error": err}, true)
		}

		defer resp2.Body.Close()

		// get the CK key on the sys.scripts page
		scriptsEndpoint, err := url.Parse(config.GetString("app.core.rest.url") + "/sys.scripts.do")

		if err != nil {
			conf.Err("The ck key could not be found from the HTML source!", log.Fields{"error": err}, true)
		}

		// get the CK key
		resp3, err := client.Get(scriptsEndpoint.String())

		if err != nil {
			conf.Err("There was an error while making the request!", log.Fields{"error": err}, true)
		}

		defer resp3.Body.Close()

		body, err = ioutil.ReadAll(resp3.Body)

		if err != nil {
			conf.Err("There was an error while reading the response body!", log.Fields{"error": err}, true)
		}

		// search for the CK key in the HTML string, present differently that at the login...
		re = regexp.MustCompile(`<input name="sysparm_ck" type="hidden" value="([a-z0-9]{72})`)
		match = re.FindStringSubmatch(string(body))

		if match == nil {
			conf.Err("There was an error while getting the ck token!", log.Fields{"error": err}, true)
		}

		ckToken = match[1]

		// build form values for script execution
		form = url.Values{}
		form.Add("script", string(dat))
		// automatically create a record for rollback
		form.Add("record_for_rollback", "on")
		form.Add("quota_managed_transaction", "on")
		form.Add("sys_scope", scopeSysID)
		form.Add("runscript", "Run script")
		form.Add("sysparm_ck", ckToken)

		resp4, err := client.PostForm(scriptsEndpoint.String(), form)

		if err != nil {
			conf.Err("There was an error while executing the script!", log.Fields{"error": err}, true)
		}

		defer resp4.Body.Close()
		// read the result
		body, err = ioutil.ReadAll(resp4.Body)

		if outputJSON, _ := rootCmd.Flags().GetBool("json"); !outputJSON {
			fmt.Println(string(body))
		} else {
			log.WithFields(log.Fields{"result": string(body)}).Info("Script execution response!")
		}
	},
}
