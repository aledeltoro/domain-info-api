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
func NewHost(URL string) *Host {

	var host Host

	servers := AddServers(URL)
	siteInfo := scraping.FetchWebsiteInfo(URL)
	status := sslAPI.SslGet(URL).Status

	host = Host{
		Servers:        servers,
		ServersChanged: false,
		Grade:          GetLowestGrade(servers),
		PreviousGrade:  "",
		Logo:           siteInfo.Logo,
		Title:          siteInfo.Title,
		IsDown:         statusMessages[status],
	}

	return &host

}
