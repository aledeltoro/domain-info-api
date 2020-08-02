package hostinfo

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	wrappedErr "domain-info-api/platform/errorhandling"
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
func NewDomain(URL string) (*Domain, *wrappedErr.Error) {

	var domainObject Domain
	var hostObject *Host

	hostObject, customErr := NewHost(URL)
	if customErr != nil {
		return &Domain{}, customErr
	}

	domainObject = Domain{
		Name:      URL,
		HostInfo:  *hostObject,
		CreatedAt: time.Now(),
	}

	return &domainObject, nil

}

// InsertDomain inserts a record into the "host" database
func (c *Connection) InsertDomain(domain *Domain) *wrappedErr.Error {

	var customErr *wrappedErr.Error

	insertDomainStmt, err := c.DB.Prepare(`
	INSERT INTO 
		host (domain_name, server_changed, ssl_grade, previous_ssl_grade, logo, title, is_down, created_at) 
	VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8) 
	RETURNING id
	`)
	if err != nil {
		errMessage := fmt.Sprintf("Invalid query statement: %s", err.Error())
		customErr = wrappedErr.New(500, "InsertDomain", errMessage)
		log.Println(customErr)
		return customErr
	}

	defer insertDomainStmt.Close()

	host := domain.HostInfo

	record := insertDomainStmt.QueryRow(domain.Name, host.ServersChanged, host.Grade, host.PreviousGrade, host.Logo, host.Title, host.IsDown, domain.CreatedAt)
	if err != nil {
		errMessage := fmt.Sprintf("Query operation failed: %s", err.Error())
		customErr = wrappedErr.New(500, "InsertDomain", errMessage)
		log.Println(customErr)
		return customErr
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
		errMessage := fmt.Sprintf("Query operation failed: %s", err.Error())
		customErr = wrappedErr.New(500, "InsertDomain", errMessage)
		log.Println(customErr)
		return customErr
	}

	defer insertServerStmt.Close()

	for i := 0; i < len(host.Servers); i++ {

		server := host.Servers[i]

		_, err := insertServerStmt.Exec(server.Address, server.SslGrade, server.Country, server.Owner, lastInsertID)
		if err != nil {
			errMessage := fmt.Sprintf("Query operation failed: %s", err.Error())
			customErr = wrappedErr.New(500, "InsertDomain", errMessage)
			log.Println(customErr)
			return customErr
		}

	}

	return nil

}

// GetAllDomains returns a slice of domains from the database
func (c *Connection) GetAllDomains() (*Items, *wrappedErr.Error) {

	var items Items
	var customErr *wrappedErr.Error

	rows, err := c.DB.Query("SELECT * FROM host")
	if err != nil {
		errMessage := fmt.Sprintf("Query operation failed: %s", err.Error())
		customErr = wrappedErr.New(500, "GetAllDomains", errMessage)
		log.Println(customErr)
		return &Items{}, customErr
	}

	defer rows.Close()

	var id int
	var serverChanged, isDown bool
	var name, grade, previousGrade, logo, title string
	var createdAt time.Time

	for rows.Next() {

		var err = rows.Scan(&id, &name, &serverChanged, &grade, &previousGrade, &logo, &title, &isDown, &createdAt)
		if err != nil {
			errMessage := fmt.Sprintf("Row scan failed: %s", err.Error())
			customErr = wrappedErr.New(500, "GetAllDomains", errMessage)
			log.Println(customErr)
			return &Items{}, customErr
		}

		servers, customErr := c.getAllServers(id)
		if customErr != nil {
			return &Items{}, customErr
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
func (c *Connection) CheckDomainExists(domainName string) (*Domain, bool, *wrappedErr.Error) {

	var customErr *wrappedErr.Error

	stmt, err := c.DB.Prepare(`
	SELECT
		host.id, host.ssl_grade, host.created_at  
	FROM 
		host 
	WHERE 
		host.domain_name=$1
	`)
	if err != nil {
		errMessage := fmt.Sprintf("Invalid query statement: %s", err.Error())
		customErr = wrappedErr.New(500, "CheckDomainExists", errMessage)
		log.Println(customErr)
		return &Domain{}, false, customErr
	}

	defer stmt.Close()

	var hostID int
	var currentGrade string
	var createdAt time.Time

	err = stmt.QueryRow(domainName).Scan(&hostID, &currentGrade, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &Domain{}, false, nil
		}
		errMessage := fmt.Sprintf("Query operation failed: %s", err.Error())
		customErr = wrappedErr.New(500, "CheckDomainExists", errMessage)
		log.Println(customErr)
		return &Domain{}, false, customErr
	}

	if diff := checkTimeDiffNow(createdAt); diff >= 1 {

		oldServers, customErr := c.getAllServers(hostID)
		if customErr != nil {
			return &Domain{}, false, customErr
		}

		newServers, customErr := AddServers(domainName)
		if customErr != nil {
			return &Domain{}, false, customErr
		}

		newGrade := GetLowestGrade(newServers)

		serverChanged := haveServersChanged(newServers, oldServers)

		if serverChanged {

			customErr = c.updateAllServers(newServers, hostID)
			if customErr != nil {
				return &Domain{}, false, customErr
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
			errMessage := fmt.Sprintf("Invalid query statement: %s", err.Error())
			customErr = wrappedErr.New(500, "CheckDomainExists", errMessage)
			log.Println(customErr)
			return &Domain{}, false, customErr
		}

		defer stmt.Close()

		_, err = stmt.Exec(serverChanged, newGrade, currentGrade, time.Now(), hostID)
		if err != nil {
			errMessage := fmt.Sprintf("Query operation failed: %s", err.Error())
			customErr = wrappedErr.New(500, "CheckDomainExists", errMessage)
			log.Println(customErr)
			return &Domain{}, false, customErr
		}

	}

	domainObject, customErr := c.GetDomain(domainName)
	if customErr != nil {
		return &Domain{}, false, customErr
	}

	return domainObject, true, nil

}

// GetDomain returns a single domain specified by the domain name
func (c *Connection) GetDomain(domainName string) (*Domain, *wrappedErr.Error) {

	var customErr *wrappedErr.Error

	stmt, err := c.DB.Prepare("SELECT * FROM host WHERE host.domain_name=$1")
	if err != nil {
		errMessage := fmt.Sprintf("Invalid query statement: %s", err.Error())
		customErr = wrappedErr.New(500, "GetDomain", errMessage)
		log.Println(customErr)
		return &Domain{}, customErr
	}

	defer stmt.Close()

	row := stmt.QueryRow(domainName)

	var id int
	var grade, previousGrade, logo, title string
	var serversChanged, isDown bool
	var createdAt time.Time

	err = row.Scan(&id, &domainName, &serversChanged, &grade, &previousGrade, &logo, &title, &isDown, &createdAt)
	if err != nil {
		errMessage := fmt.Sprintf("Row scan failed: %s", err.Error())
		customErr = wrappedErr.New(500, "GetDomain", errMessage)
		log.Println(customErr)
		return &Domain{}, customErr
	}

	servers, customErr := c.getAllServers(id)
	if customErr != nil {
		return &Domain{}, customErr
	}

	domainObject := Domain{
		Name: domainName,
		HostInfo: Host{
			Servers:        servers,
			ServersChanged: serversChanged,
			Grade:          grade,
			PreviousGrade:  previousGrade,
			Logo:           logo,
			Title:          title,
			IsDown:         isDown,
		},
		CreatedAt: createdAt,
	}

	return &domainObject, nil

}

// Returns from the database a slice of servers for a given host id
func (c *Connection) getAllServers(hostID int) ([]Server, *wrappedErr.Error) {

	var servers []Server
	var newErr *wrappedErr.Error

	stmt, err := c.DB.Prepare(`
	SELECT 
		server.address, server.ssl_grade, server.country, server.owner 
	FROM 
		server 
	WHERE 
		server.host_id=$1
	`)
	if err != nil {
		errMessage := fmt.Sprintf("Invalid query statement: %s", err.Error())
		newErr = wrappedErr.New(500, "getAllServers", errMessage)
		log.Println(newErr)
		return []Server{}, newErr
	}

	defer stmt.Close()

	rows, err := stmt.Query(hostID)
	if err != nil {
		errMessage := fmt.Sprintf("Query operation failed: %s", err.Error())
		newErr = wrappedErr.New(500, "getAllServers", errMessage)
		log.Println(newErr)
		return []Server{}, newErr
	}

	defer rows.Close()

	var address, grade, country, owner string

	for rows.Next() {

		err := rows.Scan(&address, &grade, &country, &owner)
		if err != nil {
			errMessage := fmt.Sprintf("Row scan failed: %s", err.Error())
			newErr = wrappedErr.New(500, "getAllServers", errMessage)
			log.Println(newErr)
			return []Server{}, newErr
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

func (c *Connection) updateAllServers(newServers []Server, hostID int) *wrappedErr.Error {

	var customErr *wrappedErr.Error

	deleteServerStmt, err := c.DB.Prepare(`
	DELETE FROM server
	WHERE host_id = $1;
	`)
	if err != nil {
		errMessage := fmt.Sprintf("Invalid query statement: %s", err.Error())
		customErr = wrappedErr.New(500, "updateAllServers", errMessage)
		log.Println(customErr)
		return customErr
	}

	defer deleteServerStmt.Close()

	_, err = deleteServerStmt.Exec(hostID)
	if err != nil {
		errMessage := fmt.Sprintf("Query operation failed: %s", err.Error())
		customErr = wrappedErr.New(500, "updateAllServers", errMessage)
		log.Println(customErr)
		return customErr
	}

	insertServerStmt, err := c.DB.Prepare(`
	INSERT INTO
		server (address, ssl_grade, country, owner, host_id)
	VALUES
		($1, $2, $3, $4, $5)
	`)
	if err != nil {
		errMessage := fmt.Sprintf("Invalid query statement: %s", err.Error())
		customErr = wrappedErr.New(500, "updateAllServers", errMessage)
		log.Println(customErr)
		return customErr
	}

	defer insertServerStmt.Close()

	for i := 0; i < len(newServers); i++ {

		server := newServers[i]

		_, err := insertServerStmt.Exec(server.Address, server.SslGrade, server.Country, server.Owner, hostID)
		if err != nil {
			errMessage := fmt.Sprintf("Query operation failed: %s", err.Error())
			customErr = wrappedErr.New(500, "updateAllServers", errMessage)
			log.Println(customErr)
			return customErr
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

		if oldServers[i] != newServers[i] {
			return true
		}

	}

	return false

}
