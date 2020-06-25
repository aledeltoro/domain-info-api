package hostinfo

import (
	sslAPI "domain-info-api/platform/ssllabs"
	scraping "domain-info-api/platform/webscraping"
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

// NewHost return a Host struct with about the given URL
func NewHost(URL string) (*Host, error) {

	var host Host

	servers, err := AddServers(URL)
	if err != nil {
		return &Host{}, err
	}
	
	siteInfo, err := scraping.FetchWebsiteInfo(URL)
	if err != nil {
		return &Host{}, err
	}
	
	responseObject, err := sslAPI.SslGet(URL)
	if err != nil {
		return &Host{}, err
	}

	status := responseObject.Status

	host = Host{
		Servers:        servers,
		ServersChanged: false,
		Grade:          GetLowestGrade(servers),
		PreviousGrade:  "",
		Logo:           siteInfo.Logo,
		Title:          siteInfo.Title,
		IsDown:         statusMessages[status],
	}

	return &host, nil

}
