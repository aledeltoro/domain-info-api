package webscraping

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// WebsiteInfo represents the scraped data from a given domain
type WebsiteInfo struct {
	Title string
	Logo  string
}

// FetchWebsiteInfo returns a new instance of WebsiteInfo w
func FetchWebsiteInfo(domain string) WebsiteInfo {

	var siteInfo WebsiteInfo
	
	document := scrapeDocument(domain)
	siteInfo.fetchTitle(document)
	siteInfo.fetchLogo(document)

	return siteInfo

}

func scrapeDocument(domain string) *goquery.Document {

	response, err := http.Get(domain)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return document

}


func (w *WebsiteInfo) fetchTitle(document *goquery.Document) {

	title := document.Find("title").Text()
	w.Title = title

}

func (w *WebsiteInfo) fetchLogo(document *goquery.Document) {

	document.Find("link").Each(func(index int, element *goquery.Selection) {

		rel, exists := element.Attr("rel")
		if exists {

			if rel == "shortcut icon" {
				logo, exists := element.Attr("href")
				if exists {
					w.Logo = logo
				}
			}

		}

	})

}