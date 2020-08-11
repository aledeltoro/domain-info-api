package hostinfo

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var (
	testDomain = Domain{
		Name:      "test.com",
		HostInfo:  testHost,
		CreatedAt: time.Now(),
	}

	testHost = Host{
		Servers: []Server{
			Server{
				Address:  "server1",
				SslGrade: "B",
				Country:  "US",
				Owner:    "Amazon.com, Inc.",
			},
			Server{
				Address:  "server2",
				SslGrade: "A+",
				Country:  "US",
				Owner:    "Amazon.com, Inc.",
			},
			Server{
				Address:  "server3",
				SslGrade: "A",
				Country:  "US",
				Owner:    "Amazon.com, Inc.",
			},
		},
		ServersChanged: true,
		Grade:          "B",
		PreviousGrade:  "A+",
		Logo:           "https://server.com/icon.png",
		Title:          "Title of test page",
		IsDown:         false,
	}

	mockConnection = Connection{}
)

func newMock() (*sql.DB, sqlmock.Sqlmock) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock

}

func setUpTables() (hostRows, serverRows *sqlmock.Rows) {

	hostRows = sqlmock.NewRows([]string{"id", "domain_name", "server_changed", "ssl_grade", "previous_ssl_grade", "logo", "title", "is_down", "created_at"})
	serverRows = sqlmock.NewRows([]string{"id", "address", "ssl_grade", "country", "owner", "host_id"})

	return

}

func TestInsertDomain(t *testing.T) {

	db, mock := newMock()

	insertDomainQuery := `
	INSERT INTO 
		host (domain_name, server_changed, ssl_grade, previous_ssl_grade, logo, title, is_down, created_at) 
	VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8) 
	RETURNING id
	`
	insertServerQuery := `
	INSERT INTO 
		server (address, ssl_grade, country, owner, host_id) 
	VALUES 
		($1, $2, $3, $4, $5)
	`
	hostID := 0

	domainStmt := mock.ExpectPrepare(insertDomainQuery)
	domainStmt.ExpectQuery().
		WithArgs(testDomain.Name, testHost.ServersChanged, testHost.Grade, testHost.PreviousGrade, testHost.Logo, testHost.Title, testHost.IsDown, testDomain.CreatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow(hostID))

	serverStmt := mock.ExpectPrepare(insertServerQuery)

	for i := 0; i < len(testHost.Servers); i++ {

		server := testHost.Servers[i]

		_ = serverStmt.ExpectExec().
			WithArgs(server.Address, server.SslGrade, server.Country, server.Owner, hostID).
			WillReturnResult(sqlmock.NewResult(0, 1))

	}

	mockConnection.DB = db

	customErr := mockConnection.InsertDomain(&testDomain)
	if customErr != nil {
		t.Errorf("didn't expect an error: %s", customErr)
	}

	err := mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("expectations were not met: %s", err)
	}

}

func TestGetAllDomains(t *testing.T) {

	db, mock := newMock()
	hostRows, serverRows := setUpTables()

	query := "SELECT * FROM host"

	serverQuery := `
	SELECT 
		server.address, server.ssl_grade, server.country, server.owner 
	FROM 
		server 
	WHERE 
		server.host_id=$1
	`

	for i := 0; i < 3; i++ {

		server := testHost.Servers[i]

		hostRows.AddRow(i, testDomain.Name, testHost.ServersChanged, testHost.Grade, testHost.PreviousGrade, testHost.Logo, testHost.Title, testHost.IsDown, testDomain.CreatedAt)

		serverRows.AddRow(i, server.Address, server.SslGrade, server.Country, server.Owner, i)

	}

	mock.ExpectQuery(query).WillReturnRows(hostRows)

	for i := 0; i < 3; i++ {

		serverStmt := mock.ExpectPrepare(serverQuery)
		serverStmt.ExpectQuery().
			WithArgs(i).
			WillReturnRows(sqlmock.NewRows([]string{"address", "ssl_grade", "country", "owner"}))

	}

	mockConnection.DB = db

	_, customErr := mockConnection.GetAllDomains()
	if customErr != nil {
		t.Errorf("didn't expect an error: %s", customErr)
	}

	err := mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("expectations were not met: %s", err)
	}

}

func TestGetDomain(t *testing.T) {

	db, mock := newMock()
	hostRows, serverRows := setUpTables()

	query := "SELECT * FROM host WHERE host.domain_name=$1"

	serverQuery := `
	SELECT 
		server.address, server.ssl_grade, server.country, server.owner 
	FROM 
		server 
	WHERE 
		server.host_id=$1
	`

	hostRows.AddRow(0, testDomain.Name, testHost.ServersChanged, testHost.Grade, testHost.PreviousGrade, testHost.Logo, testHost.Title, testHost.IsDown, testDomain.CreatedAt)

	for i := 0; i < 3; i++ {

		server := testHost.Servers[i]

		serverRows.AddRow(i, server.Address, server.SslGrade, server.Country, server.Owner, 0)

	}

	domainStmt := mock.ExpectPrepare(query)
	domainStmt.ExpectQuery().WithArgs("test.com").WillReturnRows(hostRows)

	serverStmt := mock.ExpectPrepare(serverQuery)
	serverStmt.ExpectQuery().WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"address", "ssl_grade", "country", "owner"}))

	mockConnection.DB = db

	_, customErr := mockConnection.getDomain("test.com")
	if customErr != nil {
		t.Errorf("didn't expect an error: %s", customErr)
	}

	err := mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("expectations were not met: %s", err)
	}

}
