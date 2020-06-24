package hostinfo

import (
	"database/sql"
	"log"
)

// Connection represents an active connection to a database
type Connection struct {
	DB *sql.DB
}

var (
	hostTableQuery = `CREATE TABLE IF NOT EXISTS host (
		id SERIAL PRIMARY KEY, 
		domain_name TEXT,
		server_changed BOOLEAN,
		ssl_grade VARCHAR(2),
		previous_ssl_grade VARCHAR(2),
		logo TEXT, 
		title TEXT,
		is_down BOOLEAN,
		created_at TIMESTAMPTZ
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

// NewConnection creates the 'host' and 'server' tables and returns a connection to the database
func NewConnection(db *sql.DB) *Connection {

	stmt, err := db.Prepare(hostTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal("Failed creation'host' table: ", err)
	}

	stmt, err = db.Prepare(serverTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal("Failed creation 'server' table: ", err)
	}

	return &Connection{
		DB: db,
	}

}
