package hostinfo

import (
	sslAPI "domain-info-api/platform/ssllabs"
	scraping "domain-info-api/platform/webscraping"
	"log"
	"time"
)

// Host represents info for a given Host
type Host struct {
	Servers        []Server `json:"servers"`
	ServersChanged bool     `json:"servers_changed"`
	Grade          string   `json:"ssl_grade"`
	PreviousGrade  string   `json:"previous_ssl_grade"`
	Logo           string   `json:"logo"`
	Title          string   `json:"title"`
	IsDown         bool     `json:"is_down"`
	CreatedAt      time.Time
}

var statusMessages = map[string]bool{
	"ERROR": true,
	"READY": false,
}

// NewHost returns a new Host d based on the given url
func NewHost(domain string) *Host {

	var host Host

	servers := AddServers(domain)
	siteInfo := scraping.FetchWebsiteInfo(domain)
	status := sslAPI.SslGet(domain).Status

	host = Host{
		Servers:        servers,
		ServersChanged: false,
		Grade:          GetLowestGrade(servers),
		PreviousGrade:  "",
		Logo:           siteInfo.Logo,
		Title:          siteInfo.Title,
		IsDown:         statusMessages[status],
		CreatedAt:      time.Now(),
	}

	return &host

}

// InsertHost inserts a record into the "host" database
func (c *Connection) InsertHost(host *Host) {

	insertHostStmt, err := c.DB.Prepare("INSERT INTO host (server_changed, ssl_grade, previous_ssl_grade, logo, title, is_down, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id")
	if err != nil {
		log.Fatal(err)
	}

	record := insertHostStmt.QueryRow(host.ServersChanged, host.Grade, host.PreviousGrade, host.Logo, host.Title, host.IsDown, host.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}

	var lastInsertID int
	record.Scan(&lastInsertID)

	insertServerStmt, err := c.DB.Prepare("INSERT INTO server (address, ssl_grade, country, owner, host_id) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(host.Servers); i++ {

		server := host.Servers[i]

		_, err := insertServerStmt.Exec(server.Address, server.SslGrade, server.Country, server.Owner, lastInsertID)
		if err != nil {
			log.Fatal(err)
		}

	}

}


