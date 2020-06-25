package ssllabs

// Response represents the response from SSL Labs API
type Response struct {
	Status    string     `json:"status"`
	EndPoints []EndPoint `json:"endpoints"`
}

// EndPoint represents info for a given server endpoint
type EndPoint struct {
	IPAddress string `json:"ipAddress"`
	Grade     string `json:"grade"`
}

var StatusMessages = map[string]bool {
	"ERROR": false,
	"READY": true,
}

