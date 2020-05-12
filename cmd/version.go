package cmd

import (
	"fmt"
	"github.com/sn-edit/sn-edit/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of sn-edit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sn-edit %s %s\n", version.GetVersion(), version.GetCommit())
	},
}
