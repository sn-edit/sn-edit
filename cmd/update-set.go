package cmd

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/sn-edit/sn-edit/cmd/updateset"
	"github.com/sn-edit/sn-edit/conf"
	"github.com/spf13/cobra"
)

var updateSetCmd = &cobra.Command{
	Use:   "updateset",
	Short: "Manage update sets for the app",
	Long: `You are able to list update sets from the instance.
Set update sets for scopes defined in the database.
Attention: An invalid scope name defaults to global scope. I warned you!`,
	Run: func(cmd *cobra.Command, args []string) {
		scopeName, err := cmd.Flags().GetString("scope")

		if err != nil {
			conf.Err("Parsing error scope flag!", log.Fields{"error": err}, true)
		}

		// list update sets for scope
		list, err := cmd.Flags().GetBool("list")

		if err != nil {
			conf.Err("Parsing error list flag!", log.Fields{"error": err}, true)
		}

		// set update sets for scope
		set, err := cmd.Flags().GetBool("set")

		if err != nil {
			conf.Err("Parsing error set flag!", log.Fields{"error": err}, true)
		}

		// truncate update sets globally
		truncate, err := cmd.Flags().GetBool("truncate")

		if err != nil {
			conf.Err("Parsing error truncate flag!", log.Fields{"error": err}, true)
		}

		// refresh flag tryncates update set data in the database
		// use this carefully
		if truncate {
			updateset.TruncateUpdateSets()
			return
		}

		if list {
			updateset.ListCommand(cmd, scopeName)
			return
		}

		if set {
			updateSetSysID, err := cmd.Flags().GetString("update_set")

			if err != nil {
				conf.Err("Parsing error update_set flag!", log.Fields{"error": err}, true)
			}

			if len(updateSetSysID) != 32 {
				log.Info("Get a list of sys_id's by calling the updateset --list command!")
				conf.Err("Please provide a valid sys_id!", log.Fields{"error": errors.New("invalid_sys_id_length")}, true)
			}

			updateset.SetCommand(scopeName, updateSetSysID)
		}
	},
}
