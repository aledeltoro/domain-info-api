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

	protocol := "https://"

	response, err := http.Get(protocol + domain)
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
