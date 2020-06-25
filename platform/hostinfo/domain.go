package hostinfo

import (
	"database/sql"
	"log"
	"time"
)

// Items represents an array of domains
type Items struct {
	Domains []Domain `json:"items"`
}

// Domain represents the host info of a given domain
type Domain struct {
	Name      string    `json:"domainName"`
	HostInfo  Host      `json:"hostInfo"`
	CreatedAt time.Time `json:"created_at"`
}

// NewDomain returns a new Domain based on the given url
func NewDomain(URL string) (*Domain, error) {

	var domainObject Domain
	var hostObject *Host

	hostObject, err := NewHost(URL)
	if err != nil {
		return &Domain{}, err
	}

	domainObject = Domain{
		Name:      URL,
		HostInfo:  *hostObject,
		CreatedAt: time.Now(),
	}

	return &domainObject, nil

}

// InsertDomain inserts a record into the "host" database
func (c *Connection) InsertDomain(domain *Domain) error {

	insertDomainStmt, err := c.DB.Prepare(`
	INSERT INTO 
		host (domain_name, server_changed, ssl_grade, previous_ssl_grade, logo, title, is_down, created_at) 
	VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8) 
	RETURNING id
	`)
	if err != nil {
		log.Println("Invalid query statement: ", err.Error())
		return err
	}

	host := domain.HostInfo

	record := insertDomainStmt.QueryRow(domain.Name, host.ServersChanged, host.Grade, host.PreviousGrade, host.Logo, host.Title, host.IsDown, domain.CreatedAt)
	if err != nil {
		log.Println("Query operation failed: ", err.Error())
		return err
	}

	var lastInsertID int
	record.Scan(&lastInsertID)

	insertServerStmt, err := c.DB.Prepare(`
	INSERT INTO 
		server (address, ssl_grade, country, owner, host_id) 
	VALUES 
		($1, $2, $3, $4, $5)
	`)
	if err != nil {
		log.Println("Invalid query statement: ", err.Error())
		return err
	}

	for i := 0; i < len(host.Servers); i++ {

		server := host.Servers[i]

		_, err := insertServerStmt.Exec(server.Address, server.SslGrade, server.Country, server.Owner, lastInsertID)
		if err != nil {
			log.Println("Query operation failed: ", err.Error())
			return err
		}

	}

	return nil

}

// GetAllDomains returns a slice of domains from the database
func (c *Connection) GetAllDomains() (*Items, error) {

	var items Items

	rows, err := c.DB.Query("SELECT * FROM host")
	if err != nil {
		log.Println("Query operation failed: ", err.Error())
		return &Items{}, err
	}

	defer rows.Close()

	var id int
	var serverChanged, isDown bool
	var name, grade, previousGrade, logo, title string
	var createdAt time.Time

	for rows.Next() {

		err := rows.Scan(&id, &name, &serverChanged, &grade, &previousGrade, &logo, &title, &isDown, &createdAt)
		if err != nil {
			log.Println("Row scan failed: ", err.Error())
			return &Items{}, err
		}

		servers, err := c.getAllServers(id)
		if err != nil {
			return &Items{}, err
		}

		domain := Domain{
			Name: name,
			HostInfo: Host{
				Servers:        servers,
				ServersChanged: serverChanged,
				Grade:          grade,
				PreviousGrade:  previousGrade,
				Logo:           logo,
				Title:          title,
				IsDown:         isDown,
			},
			CreatedAt: createdAt,
		}

		items.Domains = append(items.Domains, domain)

	}

	return &items, nil

}


// CheckDomainExists returns the given domain from the database if it already exists
func (c *Connection) CheckDomainExists(domainName string) (*Domain, bool, error) {

	stmt, err := c.DB.Prepare(`
	SELECT
		host.id, host.ssl_grade, host.created_at  
	FROM 
		host 
	WHERE 
		host.domain_name=$1
	`)
	if err != nil {
		log.Println("Invalid query statement: ", err.Error())
		return &Domain{}, false, err
	}

	var hostID int
	var currentGrade string
	var createdAt time.Time

	err = stmt.QueryRow(domainName).Scan(&hostID, &currentGrade, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &Domain{}, false, nil
		}
		log.Println("Query operation failed: ", err.Error())
		return &Domain{}, false, err
	}

	oldServers, err := c.getAllServers(hostID)
	if err != nil {
		return &Domain{}, false, err
	}

	if diff := checkTimeDiffNow(createdAt); diff >= 1 {

		newServers, err := AddServers(domainName)
		if err != nil {
			return &Domain{}, false, err
		}

		newGrade := GetLowestGrade(newServers)

		serverChanged := haveServersChanged(newServers, oldServers)

		if serverChanged {

			err = c.updateAllServers(newServers, hostID)
			if err != nil {
				return &Domain{}, false, err
			}

		}

		stmt, err := c.DB.Prepare(`
		UPDATE host
		SET server_changed = $1,
				ssl_grade = $2,
				previous_ssl_grade = $3,
				created_at = $4
		WHERE
			host.id = $5
		`)
		if err != nil {
			log.Println("Invalid query statement: ", err.Error())
			return &Domain{}, false, err
		}

		_, err = stmt.Exec(serverChanged, newGrade, currentGrade, time.Now(), hostID)
		if err != nil {
			log.Println("Query operation failed: ", err.Error())
			return &Domain{}, false, err
		}

	}

	domainObject, err := c.GetDomain(domainName)
	if err != nil {
		return &Domain{}, false, nil
	}

	return domainObject, true, nil

}

// GetDomain returns a single domain specified by the domain name
func (c *Connection) GetDomain(domainName string) (*Domain, error) {	

	stmt, err := c.DB.Prepare("SELECT * FROM host WHERE host.domain_name=$1")
	if err != nil {
		log.Println("Invalid query statement: ", err.Error())
		return &Domain{}, err
	}

	row := stmt.QueryRow(domainName)

	var id int
	var domainGrade, previousGrade, logo, title string
	var domainServerChanged, isDown bool

	err = row.Scan(&id, &domainName, &domainServerChanged, &domainGrade, &previousGrade, &logo, &title, &isDown)
	if err != nil {
		log.Println("Row scan failed: ", err.Error())
	}

	servers, err := c.getAllServers(id)
	if err != nil {
		return &Domain{}, err
	}

	domainObject := Domain{
		Name: domainName,
		HostInfo: Host{
			Servers:        servers,
			ServersChanged: domainServerChanged,
			Grade:          domainGrade,
			PreviousGrade:  previousGrade,
			Logo:           logo,
			Title:          title,
			IsDown:         isDown,
		},
	}

	return &domainObject, nil

}

// Returns from the database a slice of servers for a given host id
func (c *Connection) getAllServers(hostID int) ([]Server, error) {

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
		log.Println("Invalid query statement: ", err.Error())
		return []Server{}, err
	}

	rows, err := stmt.Query(hostID)
	if err != nil {
		log.Println("Query operation failed: ", err.Error())
		return []Server{}, err
	}

	defer rows.Close()

	var address, grade, country, owner string

	for rows.Next() {

		err := rows.Scan(&address, &grade, &country, &owner)
		if err != nil {
			log.Println("Row scan failed: ", err.Error())
			return []Server{}, err
		}

		server := Server{
			Address:  address,
			SslGrade: grade,
			Country:  country,
			Owner:    owner,
		}

		servers = append(servers, server)

	}

	return servers, nil

}


func (c *Connection) updateAllServers(newServers []Server, hostID int) error {

	stmt, err := c.DB.Prepare(`
	UPDATE server
	SET	address = $1,
			ssl_grade = $2,
			country = $3,
			owner = $4
	WHERE
		server.host_id = $5
	`)
	if err != nil {
		log.Println("Invalid query statement: ", err.Error())
		return err
	}

	for i := 0; i < len(newServers); i++ {

		server := newServers[i]

		_, err := stmt.Exec(server.Address, server.SslGrade, server.Country, server.Owner, hostID)
		if err != nil {
			log.Println("Query operation failed: ", err.Error())
			return err
		}

	}

	return nil

}

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
