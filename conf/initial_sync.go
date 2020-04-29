package db

import (
	"fmt"
	"os"
)

func BuildTables() {
	dbc := GetDB()

	sqlStmt := `
    CREATE TABLE IF NOT EXISTS entry(id integer primary key autoincrement, sys_id text, sys_name text)
    CREATE INDEX IF NOT EXISTS entries ON entry(sys_id);
    `

	_, err := dbc.Exec(sqlStmt)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
