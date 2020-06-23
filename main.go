package main

import (
	"database/sql"
	// "fmt"
	// scraping "domain-info-api/platform/webscraping"

	// sslAPI "domain-info-api/platform/ssllabs"
	// whoisAPI "domain-info-api/platform/whoisrecord"

	"log"

	_ "github.com/lib/pq"
)

var (
	connectionString = "postgres://root@localhost:26257/domain_info_api?sslmode=disable"
	hostTableQuery   = `CREATE TABLE IF NOT EXISTS host (
		id SERIAL PRIMARY KEY, 
		server_changed BOOLEAN,
		ssl_grade VARCHAR(2),
		previous_ssl_grade VARCHAR(2),
		logo TEXT, title TEXT,
		is_down BOOLEAN
	);`
	serverTableQuery = `CREATE TABLE IF NOT EXISTS server (
		id SERIAL PRIMARY KEY,
		address TEXT,
		ssl_grade VARCHAR(2),
		country CHAR(2),
		owner TEXT,
		host_id INTEGER,
		FOREIGN KEY (host_id) REFERENCES host(id)
	);`
)

func main() {

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	stmt, err := db.Prepare(hostTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal("Failed host table creation: ", err)
	}

	stmt, err = db.Prepare(serverTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal("Failed server table creation: ", err)
	}

	// responseSSL := sslAPI.SslGet("truora.com")
	// fmt.Println(responseSSL)

	// responseWhoIs := whoisAPI.WhoIsGet("34.193.69.252")
	// fmt.Println(responseWhoIs)

	// var website scraping.Website

	// document := scraping.ScrapeDocument("https://www.truora.com")

	// website.FetchTitle(document)
	// website.FetchLogo(document)

	// fmt.Println(website)

}
