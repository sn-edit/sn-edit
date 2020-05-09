package cmd

import (
	"github.com/0x111/sn-edit/cmd/updateset"
	log "github.com/sirupsen/logrus"
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
			log.WithFields(log.Fields{"error": err, "flag": "updateset.scope"}).Error("Parsing error!")
			return
		}

		// list update sets for scope
		list, err := cmd.Flags().GetBool("list")

		if err != nil {
			log.WithFields(log.Fields{"error": err, "flag": "list"}).Error("Parsing error!")
			return
		}

		// set update sets for scope
		set, err := cmd.Flags().GetBool("set")

		if err != nil {
			log.WithFields(log.Fields{"error": err, "flag": "list"}).Error("Parsing error!")
			return
		}

		if list {
			updateset.ListCommand(scopeName)
			return
		}

		if set {
			updateSetSysID, err := cmd.Flags().GetString("update_set")

			if err != nil {
				log.WithFields(log.Fields{"error": err, "flag": "updateset.scope"}).Error("Parsing error!")
				return
			}

			if len(updateSetSysID) != 32 {
				log.WithFields(log.Fields{"error": "sys_id.length", "flag": "updateset.update_set"}).Error("The update_set flag has to be a sys_id of Update Set! See --list to get the curenty avaliable entries!")
			}

			updateset.SetCommand(scopeName, updateSetSysID)
		}
	},
}
