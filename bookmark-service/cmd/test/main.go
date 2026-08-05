package main

import (
	"github.com/HemlockPham7/golang-system-design/pkg/sqldb"
)

func main() {
	dbClient, err := sqldb.NewClient("")
	if err != nil {
		panic(err)
	}

	err = sqldb.MigrateSQLDB(dbClient, "file://./migration", "steps", -1)
	if err != nil {
		panic(err)
	}
}
