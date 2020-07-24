package whoisrecord

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	wrappedErr "domain-info-api/platform/errorhandling"

	"github.com/valyala/fasthttp"
)

// Get returns the registrant information of the specified IP
func Get(IP string) (*Response, *wrappedErr.Error) {

	var customErr *wrappedErr.Error

	whoIsAPI := fmt.Sprintf("https://www.whoisxmlapi.com/whoisserver/WhoisService?apiKey=%s&outputFormat=json&domainName=", os.Getenv("WHOIS_API_KEY"))

	_, body, err := fasthttp.Get(nil, whoIsAPI+IP)
	if err != nil {
		errMessage := fmt.Sprintf("WhoisXML API consumption failed: %s", err.Error())
		customErr = wrappedErr.New(500, "Get", errMessage)
		log.Println(customErr)
		return &Response{}, customErr
	}

	var responseObject Response

	err = json.Unmarshal(body, &responseObject)
	if err != nil {
		errMessage := fmt.Sprintf("JSON enconding failed: %s", err.Error())
		customErr = wrappedErr.New(500, "Get", errMessage)
		log.Println(customErr)
		return &Response{}, customErr
	}

	return &responseObject, nil

}
