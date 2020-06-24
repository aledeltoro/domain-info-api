package ssllabs

import (
	"encoding/json"
	"log"

	"github.com/valyala/fasthttp"
)

const sslAPI = "https://api.ssllabs.com/api/v3/analyze?host="

// SslGet returns status and endpoints of the specified domain
func SslGet(domain string) *Response {

	_, body, err := fasthttp.Get(nil, sslAPI+domain)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject Response
	json.Unmarshal(body, &responseObject)

	return &responseObject

}
