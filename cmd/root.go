package cmd

import (
	"fmt"
	"github.com/mbndr/figlet4go"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/sn-edit/sn-edit/api"
	"github.com/sn-edit/sn-edit/conf"
	"github.com/sn-edit/sn-edit/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"runtime"
	"strings"
)

var (
	// Used for flags.
	cfgFile string
)

// commands list
var rootCmd = &cobra.Command{
	Use:   "sn-edit",
	Short: "An editor for developing stuff for Servicenow locally",
	Long: `sn-edit provides you a simple and easy way to edit and sync your files from your Servicenow instance
the app is lightweight and easy to use. It will give you a lot of options to work on your code locally, while syncing
to Servicenow.`,
	Version: fmt.Sprintf("%s %s %s/%s", version.GetVersion(), strings.TrimSpace(version.GetCommit()), runtime.GOOS, runtime.GOARCH),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func er(msg interface{}) {
	log.WithFields(log.Fields{"error": msg}).Error("Error")
	os.Exit(1)
}

func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()

		if err != nil {
			er(err)
		}

		// Search config in home directory with name ".sn-edit" (without extension).
		viper.AddConfigPath("./_config/")
		viper.AddConfigPath(home)
		viper.SetConfigName(".sn-edit")
	}

	viper.AutomaticEnv()

	// exclude banner if json output requested
	if outputJSON, _ := rootCmd.Flags().GetBool("json"); !outputJSON {
		if runtime.GOOS != "windows" {
			PrintBanner()
		}
	}

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	if outputJSON, _ := rootCmd.Flags().GetBool("json"); outputJSON {
		log.SetFormatter(&log.JSONFormatter{})
	}

	// do not write out text for json output
	if err := viper.ReadInConfig(); err == nil {
		log.WithFields(log.Fields{"config": viper.ConfigFileUsed()}).Info("Using config file")
	}

	conf.SetConfig(viper.GetViper())
	// Validate the config file
	conf.ValidateConfig()
	// Validate table data if correct
	conf.ValidateTableData()
	// Set the log level
	conf.SetLoggerLevel()
	// connect to db
	conf.ConnectDB()
	// setup database
	conf.BuildTables()
	// setup http client that we will use throughout the app
	api.SetupClient()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		er(err)
	}
}

func PrintBanner() {
	ascii := figlet4go.NewAsciiRender()

	// Adding the colors to RenderOptions
	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{
		figlet4go.ColorRed,
		figlet4go.ColorMagenta,
	}

	renderStr, _ := ascii.RenderOpts("sn-edit", options)
	fmt.Print(renderStr)
}

func init() {
	// load config set, also parse env variables
	cobra.OnInitialize(initConfig)
	// config file
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sn-edit.yaml)")
	// json output formatting
	rootCmd.PersistentFlags().BoolP("json", "", false, "set this if you want sn-edit to output json to stdout")
	// download command flags
	downloadEntryCmd.Flags().StringP("table", "t", "", "the table from where sn-edit should get the entry from")
	downloadEntryCmd.Flags().StringP("sys_id", "", "", "the sys_id of the entry which you would like to get")
	// upload command flags
	uploadEntryCmd.Flags().StringP("table", "t", "", "the table from where sn-edit should get the entry from")
	uploadEntryCmd.Flags().StringP("sys_id", "", "", "the sys_id of the entry which you would like to get")
	uploadEntryCmd.Flags().StringP("fields", "f", "", "provide one or more fields, comma separated (example: \"name,script,active\")")
	uploadEntryCmd.Flags().StringP("update_set", "", "", "the sys_id of an update set, you need to list the update sets before using this (example: \"<sys_id>\")")
	// update set flags
	updateSetCmd.Flags().BoolP("list", "", false, "list update sets for the scope provided")
	updateSetCmd.Flags().BoolP("set", "", false, "set update sets for the scope provided")
	updateSetCmd.Flags().BoolP("truncate", "", false, "use this to truncate the update sets and force the reload from the instance")
	updateSetCmd.Flags().StringP("scope", "", "global", "the name of the scope (example: \"global\")")
	updateSetCmd.Flags().StringP("update_set", "", "global", "the sys_id of the update_set (example: \"<sys_id>\")")
	// execute scripts flags
	executeScriptsCmd.Flags().StringP("file", "", "", "recommended use is a fullpath to the file, but you can also specify relative paths from the POV of the binary. (example: \"/home/user/background-scripts/some-script.js\")")
	executeScriptsCmd.Flags().StringP("scope", "", "global", "the name of the scope, defaults to global (example: \"global\")")
	// search flag
	searchCmd.Flags().StringP("table", "", "", "the table you want to search the entries in")
	searchCmd.Flags().StringP("fields", "", "", "comma separated list of field names, if existent will be merged with tableconfig fields for this table")
	searchCmd.Flags().StringP("encoded_query", "", "", "the encoded query we should use when searching")
	searchCmd.Flags().Int64P("limit", "", 1, "limit of the records that are returned from the API")
	//rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(downloadEntryCmd)
	rootCmd.AddCommand(uploadEntryCmd)
	rootCmd.AddCommand(updateSetCmd)
	rootCmd.AddCommand(executeScriptsCmd)
	rootCmd.AddCommand(searchCmd)
}
