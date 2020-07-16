package webscraping

import (
	"fmt"
	"testing"

	wrappedErr "domain-info-api/platform/errorhandling"

)

func TestFetchWebsiteInfo(t *testing.T) {

	assertWebsiteInfo := func(t *testing.T, got WebsiteInfo) {
		t.Helper()

		if got.Title == "" || got.Logo == "" {
			t.Errorf("expected WebsiteInfo object to be defined")
		}
	}

	assertNoError := func(t *testing.T, err *wrappedErr.Error) {
		t.Helper()

		if err != nil {
			t.Error("got an error, but didn't want one")
		}
	}

	var testDomains = []string{"github.com", "tesla.com"}

	for _, domain := range testDomains {
		message := fmt.Sprintf("%s website info", domain)
		t.Run(message, func(t *testing.T) {
			got, err := FetchWebsiteInfo(domain)
			assertWebsiteInfo(t, got)
			assertNoError(t, err)
		})
	}

}
