package extract

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/kaffeeshare/data"
)

// ogtags extracts the og:title, og:image, ... tags from a webpage
func defaultHTML(i *data.Item, sourceURL string, doc *goquery.Document) {
	fmt.Println("Running OG extract. " + sourceURL)

	selection := doc.Find("title")
	if len(selection.Nodes) != 0 {
		i.Caption = selection.Nodes[0].FirstChild.Data
	}

	selection = doc.Find("meta[property*='og']")

	for _, e := range selection.Nodes {
		m := htmlAttributeToMap(e.Attr)

		if m["property"] == "og:title" {
			i.Caption = m["content"]
		}
		if m["property"] == "og:image" {
			if !govalidator.IsRequestURL(m["content"]) {
				log.Println("Invalid url in og:image. " + sourceURL)
				continue
			}
			i.ImageURL = m["content"]
		}
		if m["property"] == "og:url" {
			if !govalidator.IsRequestURL(m["content"]) {
				log.Println("Invalid url in og:url. " + sourceURL)
				continue
			}
			i.URL = m["content"]
		}
		if m["property"] == "og:description" {
			i.Description = m["content"]
		}
	}
}
