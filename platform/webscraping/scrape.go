package webscraping

import (
	"fmt"
	"log"
	"net/http"

	wrappedErr "domain-info-api/platform/errorhandling"

	"github.com/PuerkitoBio/goquery"
)

// WebsiteInfo represents the scraped data from a given domain
type WebsiteInfo struct {
	Title string
	Logo  string
}

// FetchWebsiteInfo returns a new instance of WebsiteInfo w
func FetchWebsiteInfo(domain string) (WebsiteInfo, *wrappedErr.Error) {

	var siteInfo WebsiteInfo

	document, customErr := scrapeDocument(domain)
	if customErr != nil {
		return WebsiteInfo{}, customErr
	}

	siteInfo.fetchTitle(document)
	siteInfo.fetchLogo(document)

	return siteInfo, nil

}

func scrapeDocument(domain string) (*goquery.Document, *wrappedErr.Error) {

	var customErr *wrappedErr.Error

	protocol := "https://"

	response, err := http.Get(protocol + domain)
	if err != nil {
		errMessage := fmt.Sprintf("Error: %s", err.Error())
		customErr = wrappedErr.New(500, "scrapeDocument", errMessage)
		log.Println(customErr)
		return &goquery.Document{}, customErr
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		errMessage := fmt.Sprintf("Error: %s", err.Error())
		customErr = wrappedErr.New(500, "scrapeDocument", errMessage)
		log.Println(customErr)
		return &goquery.Document{}, customErr
	}

	return document, nil

}

func (w *WebsiteInfo) fetchTitle(document *goquery.Document) {

	title := document.Find("title").Text()
	w.Title = title

}

func (w *WebsiteInfo) fetchLogo(document *goquery.Document) {

	document.Find("link").EachWithBreak(func(index int, element *goquery.Selection) bool {

		rel, exists := element.Attr("rel")

		var isIconReal bool = rel == "shortcut icon" || rel == "icon"

		if exists && isIconReal {

			logo, exists := element.Attr("href")
			if exists {
				w.Logo = logo
				return false
			}

		}

		return true

	})

}
