package updateset

import (
	log "github.com/sirupsen/logrus"
	"github.com/sn-edit/sn-edit/conf"
	"os"
)

func TruncateUpdateSets() {
	dbc := conf.GetDB()
	log.Debug("Truncating the update_set data!")

	query := "DELETE FROM update_set;"

	_, err := dbc.Exec(query)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not truncate the update_set table!")
		os.Exit(1)
	}
}
