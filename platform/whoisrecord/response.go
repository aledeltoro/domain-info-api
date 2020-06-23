package whoisrecord

// Response represents the response from the WHOIS API
type Response struct {
	WhoIsRecord RegistryData `json:"WhoisRecord"`
}

// RegistryData represents the registry data of the given domain
type RegistryData struct {
	Registry	Registrant	`json:"registryData"`
}

// Registrant represents the registrant info of a given IP address
type Registrant struct {
	RegistrantInfo Info `json:"registrant"`
}

// Info represents the information of the registrant
type Info struct {
	Organization string	`json:"organization"`
	CountryCode string	`json:"countryCode"`
}