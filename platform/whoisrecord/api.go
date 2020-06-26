package whoisrecord

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/valyala/fasthttp"
)


// WhoIsGet returns the registrant information of the specified IP
func WhoIsGet(IP string) (*Response, error) {
	
	whoIsAPI := fmt.Sprintf("https://www.whoisxmlapi.com/whoisserver/WhoisService?apiKey=%s&outputFormat=json&domainName=", os.Getenv("WHOIS_API_KEY"))

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
