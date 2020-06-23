package hostinfo

import (
	sslAPI "domain-info-api/platform/ssllabs"
	scraping "domain-info-api/platform/webscraping"
	"time"
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
	CreatedAt      time.Time
}

var statusMessage = map[string]bool{
	"ERROR": true,
	"READY": false,
}

// CreateHost returns a new Host d based on the given url
func CreateHost(domain string) *Host {

	var host Host

	servers := AddServers(domain)
	siteInfo := scraping.FetchWebsiteInfo(domain)
	status := sslAPI.SslGet(domain).Status

	host = Host{
		ServersChanged: false,
		Grade:          GetLowestGrade(servers),
		PreviousGrade:  "",
		Logo:           siteInfo.Logo,
		Title:          siteInfo.Title,
		IsDown:         statusMessage[status],
		CreatedAt:      time.Now(),
	}

	return &host

}
