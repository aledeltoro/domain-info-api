package whoisrecord

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	apiKey   = "at_y0Q4RvF1XeSEl1qvQOooAGeihm6vg"
	whoIsAPI = fmt.Sprintf("https://www.whoisxmlapi.com/whoisserver/WhoisService?apiKey=%s&outputFormat=json&domainName=", apiKey)
)

// WhoIsGet returns the registrant information of the specified IP
func WhoIsGet(IP string) *Response  {

	response, err := http.Get(whoIsAPI + IP)
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
