package hostinfo

import (
	wrappedErr "domain-info-api/platform/errorhandling"
	sslAPI "domain-info-api/platform/ssllabs"
	scraping "domain-info-api/platform/webscraping"
	"strings"
)

// Host represents info for a given Host
type Host struct {
	Servers        []Server `json:"servers"`
	ServersChanged bool     `json:"servers_changed"`
	Grade          string   `json:"ssl_grade"`
	PreviousGrade  string   `json:"previous_ssl_grade"`
	Logo           string   `json:"logo"`
	Title          string   `json:"title"`
	IsDown         bool     `json:"is_down"`
}

var statusMessages = map[string]bool{
	"ERROR": true,
	"READY": false,
}

// newHost return a Host struct with about the given URL
func newHost(URL string) (*Host, *wrappedErr.Error) {

	var host Host

	servers, customErr := addServers(URL)
	if customErr != nil {
		return &Host{}, customErr
	}

	siteInfo, customErr := scraping.FetchWebsiteInfo(URL)
	if customErr != nil {
		return &Host{}, customErr
	}

	responseObject, customErr := sslAPI.Get(URL)
	if customErr != nil {
		return &Host{}, customErr
	}

	status := responseObject.Status

	host = Host{
		Servers:        servers,
		ServersChanged: false,
		Grade:          getLowestGrade(servers),
		PreviousGrade:  "",
		Logo:           siteInfo.Logo,
		Title:          strings.TrimSpace(siteInfo.Title),
		IsDown:         statusMessages[status],
	}

	return &host, nil

}
