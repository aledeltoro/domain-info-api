package ssllabs

import (
	"encoding/json"
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

const sslAPI = "https://api.ssllabs.com/api/v3/analyze"

// SslGet returns status and endpoints of the specified domain
func SslGet(domain string) (*Response, error) {

	hostQuery := "?host="

	var responseObject Response
	var pendingResponse = true

	for pendingResponse {

		_, body, err := fasthttp.Get(nil, sslAPI+hostQuery+domain)
		if err != nil {
			log.Println("SSL API consumption failed: ", err.Error())
			return &Response{}, err
		}

		err = json.Unmarshal(body, &responseObject)
		if err != nil {
			log.Println("JSON encoding failed: ", err.Error())
			return &Response{}, err
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
