package main

import (
	"database/sql"
	hostinfo "domain-info-api/platform/hostinfo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

	router.POST("/domains", func(ctx *fasthttp.RequestCtx) {

		hostArg := ctx.URI().QueryArgs().Peek("host")

		domain := hostinfo.NewDomain(string(hostArg))
		host.InsertDomain(domain)

		ctx.Response.SetStatusCode(http.StatusCreated)
		ctx.Response.Header.SetContentType("application/json")

		json.NewEncoder(ctx).Encode(domain.HostInfo)

	})
	
	router.GET("/domains", func(ctx *fasthttp.RequestCtx) {
		
		domains := host.GetAllDomains()
		
		ctx.Response.Header.SetContentType("application/json")
		ctx.Response.SetStatusCode(http.StatusOK)
		
		err := json.NewEncoder(ctx).Encode(domains)
		if err != nil {
			log.Fatal(err)
		}
		
	})

	fmt.Println("Listening on port 3000")

	fasthttp.ListenAndServe(":3000", router.Handler)

}
