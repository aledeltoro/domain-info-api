package hostinfo

import (
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

	rowsHosts, err := c.DB.Query("SELECT * FROM host")
	if err != nil {
		log.Fatal(err)
	}

	var id int
	var serverChanged, isDown bool
	var name, grade, previousGrade, logo, title string
	var createdAt time.Time

	for rowsHosts.Next() {

		err := rowsHosts.Scan(&id, &name, &serverChanged, &grade, &previousGrade, &logo, &title, &isDown, &createdAt)
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

	var address, grade, country, owner string

	for rows.Next() {

		err := rows.Scan(&address, &grade, &country, &owner)
		if err != nil {
			log.Fatal(err)
		}

		server := Server{
			Address: address,
			SslGrade: grade,
			Country: country, 
			Owner: owner,
		}

		servers = append(servers, server)

	}

	return servers

}
