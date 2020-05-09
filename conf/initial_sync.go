package conf

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func BuildTables() {
	dbc := GetDB()
	config := GetConfig()

	sqlStmt := `
    CREATE TABLE IF NOT EXISTS entry(id integer primary key autoincrement, sys_id text, sys_name text, unique_key text, entry_table integer, sys_scope integer, last_modified integer, FOREIGN KEY(entry_table) REFERENCES entry_table(id), FOREIGN KEY(sys_scope) REFERENCES entry_scope(id));
    CREATE TABLE IF NOT EXISTS entry_table(id integer primary key autoincrement, sys_id text, name text, sys_scope integer, FOREIGN KEY(sys_scope) REFERENCES entry_scope(id));
    CREATE TABLE IF NOT EXISTS entry_scope(id integer primary key autoincrement, sys_id text, name text, update_set integer);
    CREATE TABLE IF NOT EXISTS update_set(id integer primary key autoincrement, sys_id text, name text, sys_scope integer, FOREIGN KEY(sys_scope) REFERENCES entry_scope(id));
    CREATE INDEX IF NOT EXISTS idx_entries_ids ON entry(sys_id);
    CREATE INDEX IF NOT EXISTS idx_entries_tables ON entry_table(sys_id);
    CREATE INDEX IF NOT EXISTS idx_entries_scopes ON entry_scope(sys_id);
    CREATE INDEX IF NOT EXISTS idx_update_sets ON update_set(sys_id);
    `

	if !config.GetBool("app.core.db.initialised") {
		_, err := dbc.Exec(sqlStmt)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		config.Set("app.core.db.initialised", true)
		err = config.WriteConfig()

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("There was a problem while rewriting the config file! Check the permissions please!")
			os.Exit(1)
		}
	}
}
