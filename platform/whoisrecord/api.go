package whoisrecord

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

var (
	apiKey   = "at_y0Q4RvF1XeSEl1qvQOooAGeihm6vg"
	whoIsAPI = fmt.Sprintf("https://www.whoisxmlapi.com/whoisserver/WhoisService?apiKey=%s&outputFormat=json&domainName=", apiKey)
)

// WhoIsGet returns the registrant information of the specified IP
func WhoIsGet(IP string) (*Response, error) {

	_, body, err := fasthttp.Get(nil, whoIsAPI+IP)
	if err != nil {
		log.Println("WhoisXML API consumption failed: ", err.Error())
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
