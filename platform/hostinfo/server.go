package hostinfo

import (
	wrappedErr "domain-info-api/platform/errorhandling"
	sslAPI "domain-info-api/platform/ssllabs"
	whoisAPI "domain-info-api/platform/whoisrecord"
)

// Server represents info for specific server in a given domain
type Server struct {
	Address  string `json:"address"`
	SslGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}

var grades = map[string]int{
	"A+": 7,
	"A":  6,
	"B":  5,
	"C":  4,
	"D":  3,
	"E":  2,
	"F":  1,
}

// addServers returns a slice with all of the servers of a given domain
func addServers(domain string) ([]Server, *wrappedErr.Error) {

	var servers []Server

	hostSSLData, customErr := sslAPI.Get(domain)
	if customErr != nil {
		return []Server{}, customErr
	}

	var IPAddress string

	for i := 0; i < len(hostSSLData.EndPoints); i++ {

		IPAddress = hostSSLData.EndPoints[i].IPAddress
		serverRegistry, customErr := whoisAPI.Get(IPAddress)
		if customErr != nil {
			return []Server{}, customErr
		}

		countryCode, organization := assignRegistryData(serverRegistry)

		var server = Server{
			Address:  IPAddress,
			SslGrade: hostSSLData.EndPoints[i].Grade,
			Country:  countryCode,
			Owner:    organization,
		}

		servers = append(servers, server)

	}

	return servers, nil

}

// getLowestGrade returns the lowest grade string from the array of servers
func getLowestGrade(servers []Server) string {

	var lowestGrade, currentLetter string

	lowestGrade = "A+"

	for i := 0; i < len(servers); i++ {

		currentLetter = servers[i].SslGrade

		if grades[lowestGrade] > grades[currentLetter] {
			lowestGrade = currentLetter
		}

	}

	return lowestGrade

}

func assignRegistryData(server *whoisAPI.Response) (country, owner string) {

	country = server.WhoIsRecord.Registry.RegistrantInfo.CountryCode
	owner = server.WhoIsRecord.Registry.RegistrantInfo.Organization

	if len(country) == 0 && len(owner) == 0 {

		country = server.WhoIsRecord.SubRecord[0].RegistrantInfo.CountryCode
		owner = server.WhoIsRecord.SubRecord[0].RegistrantInfo.Organization

	}

	return
}
