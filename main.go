package main

import (
	"os"
	"database/sql"
	"fmt"
	"log"

	handler "domain-info-api/handler"
	hostinfo "domain-info-api/platform/hostinfo"

	"github.com/buaazp/fasthttprouter"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/valyala/fasthttp"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
}

func main() {

	db, err := sql.Open("postgres", os.Getenv("CONNECTION_STRING"))
	if err != nil {
		log.Fatal("Failed to establish connection to database: ", err.Error())
	}

	defer db.Close()

	host, err := hostinfo.NewConnection(db)
	if err != nil {
		log.Fatal("Failed to create database's tables: ", err.Error())
	}

	router := fasthttprouter.New()

	router.POST("/domains", handler.DomainPOST(host))
	router.GET("/domains", handler.DomainGET(host))

	fmt.Println("Listening on port 3000")

	fasthttp.ListenAndServe(":3000", router.Handler)

}
