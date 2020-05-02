package cmd

import (
	"fmt"
	"github.com/0x111/sn-edit/api"
	"github.com/0x111/sn-edit/conf"
	"github.com/mbndr/figlet4go"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
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
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
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

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	conf.SetConfig(viper.GetViper())
	conf.SetLoggerLevel()
	// connect to db
	conf.ConnectDB()
	// setup database
	conf.BuildTables()
	// setup http client that we will use throughout the app
	api.SetupClient()
	// setup directory structure
	conf.SetupDirectoryStructure()
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
	// download command flags
	downloadEntryCmd.Flags().StringP("table", "", "", "the table from where sn-edit should get the entry from")
	downloadEntryCmd.Flags().StringP("sys_id", "", "", "the sys_id of the entry which you would like to get")

	uploadEntryCmd.Flags().StringP("table", "", "", "the table from where sn-edit should get the entry from")
	uploadEntryCmd.Flags().StringP("sys_id", "", "", "the sys_id of the entry which you would like to get")
	uploadEntryCmd.Flags().StringP("fields", "", "", "provide one or more fields, comma separated (example: \"name,script,active\")")
	// assign the commands to the cli
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(downloadEntryCmd)
	rootCmd.AddCommand(uploadEntryCmd)
	PrintBanner()
}
