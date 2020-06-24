package main

import (
	"database/sql"
	"fmt"
	"log"

	handler "domain-info-api/handler"
	hostinfo "domain-info-api/platform/hostinfo"

	"github.com/buaazp/fasthttprouter"
	_ "github.com/lib/pq"
	"github.com/valyala/fasthttp"
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

	host := hostinfo.NewConnection(db)

	router := fasthttprouter.New()

	router.POST("/domains", handler.DomainPOST(host))
	router.GET("/domains", handler.DomainGET(host))

	fmt.Println("Listening on port 3000")

	fasthttp.ListenAndServe(":3000", router.Handler)

}
