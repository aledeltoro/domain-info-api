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
		log.Fatal("Error loading .env file")
	}

}

func main() {

	db, err := sql.Open("postgres", os.Getenv("CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("Failed to establish connection to database: %s", err.Error())
	}

	defer db.Close()

	host, customErr := hostinfo.NewConnection(db)
	if customErr != nil {
		log.Fatal(customErr)
	}

	router := fasthttprouter.New()
	app := handler.APP{host}

	router.POST("/domains", app.DomainPOST)
	router.GET("/domains", app.DomainGET)

	fmt.Println("Listening on port 3000")

	err = fasthttp.ListenAndServe(":3000", router.Handler)
	if err != nil {
		log.Fatalf("Failed to listen to port 3000: %s", err.Error())
	}

}
