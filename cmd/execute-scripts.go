package cmd

import (
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
			log.WithFields(log.Fields{"error": err}).Error("Parsing error!")
			return
		}

		if len(scriptFile) == 0 {
			log.WithFields(log.Fields{"error": err}).Error("Please provide a file!")
			return
		}

		scopeName, err := cmd.Flags().GetString("scope")

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Parsing error!")
			return
		}

		scopeSysID, err := db.RequestScopeDataFromInstance(scopeName)

		if err != nil {
			return
		}

		if exists := file.Exists(scriptFile); exists == false {
			log.WithFields(log.Fields{"error": err, "file": scriptFile}).Error("The file could not be found! Please check the permissions and the path!")
			return
		}

		// read file contents
		dat, err := ioutil.ReadFile(scriptFile)

		//fmt.Println(string(dat))
		cookieJar, err := cookiejar.New(nil)

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Error initializing cookie jar!")
			return
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
			log.WithFields(log.Fields{"error": err}).Error("There was an error while parsing the url!")
			return
		}

		// get the CK key for login
		resp1, err := client.Get(loginUrl.String())

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("There was an error while making the request!")
			return
		}

		defer resp1.Body.Close()

		body, err := ioutil.ReadAll(resp1.Body)

		// search for the CK key in the HTML string
		re := regexp.MustCompile(`<input name="sysparm_ck" id="sysparm_ck" type="hidden" value="([a-z0-9]{72})`)
		match := re.FindStringSubmatch(string(body))

		if match == nil {
			log.WithFields(log.Fields{"error": match}).Error("There was an error while getting the ck token!")
			return
		}

		ckToken := match[1]

		form := url.Values{}
		form.Add("user_name", username)
		form.Add("user_password", passwordCredential)
		form.Add("sys_action", "sysverb_login")
		form.Add("sysparm_ck", ckToken)

		// login to the instance
		resp2, err := client.PostForm(loginUrl.String(), form)

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("There was an error while making the request!")
			return
		}

		defer resp2.Body.Close()

		// get the CK key on the sys.scripts page
		scriptsEndpoint, err := url.Parse(config.GetString("app.core.rest.url") + "/sys.scripts.do")

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("There was an error while parsing the url!")
			return
		}

		// get the CK key
		resp3, err := client.Get(scriptsEndpoint.String())

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("There was an error while making the request!")
			return
		}

		defer resp3.Body.Close()

		body, err = ioutil.ReadAll(resp3.Body)

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Could not read response body!")
			return
		}

		// search for the CK key in the HTML string, present differently that at the login...
		re = regexp.MustCompile(`<input name="sysparm_ck" type="hidden" value="([a-z0-9]{72})`)
		match = re.FindStringSubmatch(string(body))

		if match == nil {
			log.WithFields(log.Fields{"error": err}).Error("There was an error while getting the ck token!")
			return
		}

		ckToken = match[1]

		form = url.Values{}
		form.Add("script", string(dat))
		form.Add("record_for_rollback", "on")
		form.Add("quota_managed_transaction", "on")
		form.Add("sys_scope", scopeSysID)
		form.Add("runscript", "Run script")
		form.Add("sysparm_ck", ckToken)

		resp4, err := client.PostForm(scriptsEndpoint.String(), form)

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("There was an error while making the request!")
			return
		}

		defer resp4.Body.Close()
		// read the result
		body, err = ioutil.ReadAll(resp4.Body)
		fmt.Println("SCOPE: ", scopeSysID)
		if outputJSON, _ := rootCmd.Flags().GetBool("json"); !outputJSON {
			fmt.Println(string(body))
		} else {
			log.WithFields(log.Fields{"result": string(body)}).Info("Script execution response!")
		}
	},
}
