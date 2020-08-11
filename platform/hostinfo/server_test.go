package hostinfo

import (
	"fmt"
	"testing"
)

func TestGetLowestGrade(t *testing.T) {

	assertLowestGrade := func(t *testing.T, got, want string) {
		t.Helper()

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	}

	var tests = []struct {
		servers []Server
		want    string
	}{
		{servers: []Server{
			Server{
				Address:  "server1",
				SslGrade: "B",
				Country:  "US",
				Owner:    "Amazon.com, Inc.",
			},
			Server{
				Address:  "server2",
				SslGrade: "A+",
				Country:  "US",
				Owner:    "Amazon.com, Inc.",
			},
			Server{
				Address:  "server3",
				SslGrade: "A",
				Country:  "US",
				Owner:    "Amazon.com, Inc.",
			},
		}, want: "B"},
		{servers: []Server{
			Server{
				Address:  "server1",
				SslGrade: "A",
				Country:  "US",
				Owner:    "Cloudflare, Inc.",
			},
			Server{
				Address:  "server2",
				SslGrade: "A",
				Country:  "US",
				Owner:    "Cloudflare, Inc.",
			},
			Server{
				Address:  "server3",
				SslGrade: "A",
				Country:  "US",
				Owner:    "Cloudflare, Inc.",
			},
		}, want: "A"},
	}

	for _, test := range tests {
		message := fmt.Sprintf("expected SSL grade: %s", test.want)
		t.Run(message, func(t *testing.T) {
			got := getLowestGrade(test.servers)
			assertLowestGrade(t, got, test.want)
		})
	}

}
