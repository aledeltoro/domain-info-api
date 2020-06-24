package hostinfo

import (
	"database/sql"
	"log"
	"time"
)

// Domain represents the host info of a given domain
type Domain struct {
	Name     string `json:"domainName"`
	HostInfo Host   `json:"hostInfo"`
}

// NewDomain returns a new Domain based on the given url
func NewDomain(URL string) *Domain {

	var domainObject Domain

	domainObject = Domain{
		Name:     URL,
		HostInfo: *NewHost(URL),
	}

	return &domainObject

}

// InsertDomain inserts a record into the "host" database
func (c *Connection) InsertDomain(domain *Domain) {

	insertDomainStmt, err := c.DB.Prepare("INSERT INTO host (domain_name, server_changed, ssl_grade, previous_ssl_grade, logo, title, is_down, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id")
	if err != nil {
		log.Fatal(err)
	}

	host := domain.HostInfo

	record := insertDomainStmt.QueryRow(domain.Name, host.ServersChanged, host.Grade, host.PreviousGrade, host.Logo, host.Title, host.IsDown, host.CreatedAt)
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

// GetAllDomains returns a slice of domains from the database
func (c *Connection) GetAllDomains() []Domain {

	var domains []Domain

	rows, err := c.DB.Query("SELECT * FROM host")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var id int
	var serverChanged, isDown bool
	var name, grade, previousGrade, logo, title string
	var createdAt time.Time

	for rows.Next() {

		err := rows.Scan(&id, &name, &serverChanged, &grade, &previousGrade, &logo, &title, &isDown, &createdAt)
		if err != nil {
			log.Fatal(err)
		}

		domain := Domain{
			Name: name,
			HostInfo: Host{
				Servers:        c.getAllServers(id),
				ServersChanged: serverChanged,
				Grade:          grade,
				PreviousGrade:  previousGrade,
				Logo:           logo,
				Title:          title,
				IsDown:         isDown,
				CreatedAt:      createdAt,
			},
		}

		domains = append(domains, domain)

	}

	return domains

}

func (c *Connection) getAllServers(hostID int) []Server {

	var servers []Server

	stmt, err := c.DB.Prepare(`
	SELECT 
		server.address, server.ssl_grade, server.country, server.owner 
	FROM 
		server 
	WHERE 
		server.host_id=$1
	`)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := stmt.Query(hostID)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var address, grade, country, owner string

	for rows.Next() {

		err := rows.Scan(&address, &grade, &country, &owner)
		if err != nil {
			log.Fatal(err)
		}

		server := Server{
			Address:  address,
			SslGrade: grade,
			Country:  country,
			Owner:    owner,
		}

		servers = append(servers, server)

	}

	return servers

}

// CheckDomainExists returns the given domain from the database if it already exists
func (c *Connection) CheckDomainExists(domainName string) (*Domain, bool) {

	stmt, err := c.DB.Prepare(`
	SELECT
		host.id, host.ssl_grade, host.created_at  
	FROM 
		host 
	WHERE 
		host.domain_name=$1
	`)

	var id, currentGrade string
	var createdAt time.Time

	err = stmt.QueryRow(domainName).Scan(&id, &currentGrade, &createdAt)
	if err == sql.ErrNoRows {
		return &Domain{}, false
	}

	var domainObject Domain

	if diff := checkTimeDiffNow(createdAt); diff >= 1 {

		newServers := AddServers(domainName)
		newGrade := GetLowestGrade(newServers)

		stmt, err = c.DB.Prepare(`
		SELECT 
			server.address, server.ssl_grade, server.country, server.owner
		FROM
			server
		WHERE
			server.host_id=$1
		`)
		if err != nil {
			log.Fatal(err)
		}

		rows, err := stmt.Query(id)
		if err != nil {
			log.Fatal(err)
		}

		var address, serverGrade, country, owner string
		var oldServers []Server

		for rows.Next() {

			if err := rows.Scan(&address, &serverGrade, &country, &owner); err != nil {
				log.Fatal(err)
			}

			server := Server{
				Address:  address,
				SslGrade: serverGrade,
				Country:  country,
				Owner:    owner,
			}

			oldServers = append(oldServers, server)

		}

		serverChanged := haveServersChanged(newServers, oldServers)

		domainServers := oldServers

		if serverChanged {

			domainServers = newServers

			stmt, err = c.DB.Prepare(`
			UPDATE server
			SET	address = $1,
					ssl_grade = $2,
					country = $3,
					owner = $4
			WHERE
				server.host_id = $5
			`)
			if err != nil {
				log.Fatal(err)
			}

			for i := 0; i < len(newServers); i++ {

				server := newServers[i]

				_, err := stmt.Exec(server.Address, server.SslGrade, server.Country, server.Owner, id)
				if err != nil {
					log.Fatal(err)
				}

			}

		}

		stmt, err := c.DB.Prepare(`
		UPDATE host
		SET server_changed = $1,
				ssl_grade = $2,
				previous_ssl_grade = $3
		WHERE
			host.id = $4
		RETURNING 
			host.domain_name, host.server_changed, host.ssl_grade, host.previous_ssl_grade, host.logo, host.title, host.is_down
		`)
		if err != nil {
			log.Fatal(err)
		}

		row := stmt.QueryRow(serverChanged, newGrade, currentGrade, id)

		var domainName, domainGrade, previousGrade, logo, title string
		var domainServerChanged, isDown bool

		row.Scan(&domainName, &domainServerChanged, &domainGrade, &previousGrade, &logo, &title, &isDown)

		domainObject = Domain{
			Name: domainName, 
			HostInfo: Host{
				Servers: domainServers,
				ServersChanged: domainServerChanged,
				Grade: domainGrade,
				PreviousGrade: previousGrade,
				Logo: logo, 
				Title: title,
				IsDown: isDown,
			},
		}

	}

	return &domainObject, true

}

// // GetDomain returns a single domain specified by the domain name
// func (c *Connection) GetDomain(domainName string) *Domain {

// }

func checkTimeDiffNow(createdAt time.Time) float64 {

	now := time.Now()
	diff := now.Sub(createdAt).Hours()
	return diff

}

func haveServersChanged(newServers, oldServers []Server) bool {

	if len(newServers) != len(oldServers) {
		return true
	}

	for i := 0; i < len(oldServers); i++ {

		if oldServers[i].Address != newServers[i].Address {
			return true
		}

	}

	return false

}
