package whoisrecord

// Response represents the response from the WHOIS API
type Response struct {
	WhoIsRecord RegistryData `json:"WhoisRecord"`
}

// RegistryData represents the registry data of the given domain
type RegistryData struct {
	Registry  Registrant       `json:"registryData"`
	SubRecord [1]SubRegistrant `json:"subRecords"`
}

// Registrant represents the registrant info of a given IP address
type Registrant struct {
	RegistrantInfo Info `json:"registrant"`
}

// SubRegistrant represents the registrant info from an alternate response of the given IP address
type SubRegistrant struct {
	RegistrantInfo SubInfo `json:"registrant"`
}

// Info represents the information of the registrant
type Info struct {
	Organization string `json:"organization"`
	CountryCode  string `json:"countryCode"`
}

// SubInfo represents the information of the registrant from an alternate response
type SubInfo struct {
	Organization string `json:"organization"`
	CountryCode  string `json:"country"`
}
