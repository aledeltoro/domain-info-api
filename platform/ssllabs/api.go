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

		log.Printf("SSL API Status: %s", status)

		if status == "DNS" || status == "IN_PROGRESS" {
			time.Sleep(15 * time.Second)
		} else {
			pendingResponse = false
		}

	}

	return &responseObject, nil

}
