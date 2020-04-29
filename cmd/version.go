package cmd

import (
	"fmt"
	"github.com/0x111/sn-edit/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of sn-edit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sn-edit %s\n", version.GetVersion())
	},
}
