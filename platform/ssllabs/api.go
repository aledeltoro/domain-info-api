package ssllabs

import (
	"encoding/json"
	"log"

	"github.com/valyala/fasthttp"
)

const sslAPI = "https://api.ssllabs.com/api/v3/analyze"

// SslGet returns status and endpoints of the specified domain
func SslGet(domain string) (*Response, error) {

	hostQuery := "?host="

	_, body, err := fasthttp.Get(nil, sslAPI+hostQuery+domain)
	if err != nil {
		log.Println("SSL API consumption failed: ", err.Error())
		return &Response{}, err
	}

	var responseObject Response
	
	err = json.Unmarshal(body, &responseObject)
	if err != nil {
		log.Println("JSON encoding failed: ", err.Error())
		return &Response{}, err
	}

	return &responseObject, nil

}
