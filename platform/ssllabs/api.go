package ssllabs

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	wrappedErr "domain-info-api/platform/errorhandling"

	"github.com/valyala/fasthttp"
)

const sslAPI = "https://api.ssllabs.com/api/v3/analyze"

// SslGet returns status and endpoints of the specified domain
func SslGet(domain string) (*Response, *wrappedErr.Error) {

	hostQuery := "?host="

	var responseObject Response
	var pendingResponse = true
	var customErr *wrappedErr.Error
	
	startTime := time.Now()

	for pendingResponse {

		_, body, err := fasthttp.Get(nil, sslAPI+hostQuery+domain)
		if err != nil {
			errMessage := fmt.Sprintf("SSL API consumption failed: %s", err.Error())
			customErr = wrappedErr.New(500, "SslGet", errMessage)
			log.Println(customErr)
			return &Response{}, customErr
		}

		err = json.Unmarshal(body, &responseObject)
		if err != nil {
			errMessage := fmt.Sprintf("JSON encoding failed: %s", err.Error())
			customErr = wrappedErr.New(500, "SslGet", errMessage)
			log.Println(customErr)
			return &Response{}, customErr
		}

		status := responseObject.Status

		log.Printf("Domain: '%s'. SSL API Status: %s", domain, status)

		timeout := time.Now().Sub(startTime).Minutes()

		if timeout >= float64(2) {
			errMessage := fmt.Sprint("Domain could not be resolved in time. Try again later")
			customErr = wrappedErr.New(408, "SslGet", errMessage)
			log.Println(customErr)
			return &Response{}, customErr
		}


		if status == "DNS" || status == "IN_PROGRESS" {
			time.Sleep(15 * time.Second)
		} else if status == "READY" || status == "ERROR" {
			pendingResponse = false
		} else {
			errMessage := fmt.Sprintf("Unknown status found on SSL Labs API response. Try again later")
			customErr = wrappedErr.New(501, "SslGet", errMessage)
			log.Println(customErr)
			return &Response{}, customErr
		}

	}

	return &responseObject, nil

}
