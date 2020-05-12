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
	cfgFile     string
	userLicense string
)

// commands list
var rootCmd = &cobra.Command{
	Use:   "sn-edit",
	Short: "An editor for developing stuff for Servicenow locally",
	Long: `sn-edit provides you a simple and easy way to edit and sync your files from your Servicenow instance
the app is lightweight and easy to use. It will give you a lot of options to work on your code locally, while syncing
to Servicenow.`,
	Version: fmt.Sprintf("%s %s %s/%s\n", version.GetVersion(), strings.TrimSpace(version.GetCommit()), runtime.GOOS, runtime.GOARCH),
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
		viper.AddConfigPath(home)
		viper.AddConfigPath("./_config/")
		viper.SetConfigName(".sn-edit")
	}

	viper.AutomaticEnv()

	// exclude banner if json output requested
	if outputJSON, _ := rootCmd.Flags().GetBool("json"); !outputJSON {
		PrintBanner()
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
	updateSetCmd.Flags().BoolP("list", "", false, "use this to list update sets for the scope provided")
	updateSetCmd.Flags().BoolP("set", "", false, "use this to set update sets for the scope provided")
	updateSetCmd.Flags().StringP("scope", "", "global", "the name of the scope (example: \"global\")")
	updateSetCmd.Flags().StringP("update_set", "", "global", "the sys_id of the update_set (example: \"<sys_id>\")")
	//rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(downloadEntryCmd)
	rootCmd.AddCommand(uploadEntryCmd)
	rootCmd.AddCommand(updateSetCmd)
}
