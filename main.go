package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

	host, customErr := hostinfo.NewConnection(db)
	if customErr != nil {
		log.Fatal(customErr)
	}

	router := fasthttprouter.New()

	router.POST("/domains", handler.DomainPOST(host))
	router.GET("/domains", handler.DomainGET(host))

	fmt.Println("Listening on port 3000")

	err = fasthttp.ListenAndServe(":3000", router.Handler)
	if err != nil {
		log.Fatal("Failed to listen to port 3000")
	}

}
