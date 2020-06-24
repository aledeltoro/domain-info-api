package main

import (
	"database/sql"
	// hostinfo "domain-info-api/platform/hostinfo"
	"log"

	_ "github.com/lib/pq"
)

var (
	connectionString = "postgres://root@localhost:26257/domain_info_api?sslmode=disable"
)

func main() {

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// host := hostinfo.NewConnection(db)

}
