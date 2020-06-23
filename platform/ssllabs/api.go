package ssllabs

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const sslAPI = "https://api.ssllabs.com/api/v3/analyze?host="

// SslGet returns status and endpoints of the specified domain
func SslGet(domain string) *Response {

	response, err := http.Get(sslAPI + domain)
	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)
	
	return &responseObject

}
