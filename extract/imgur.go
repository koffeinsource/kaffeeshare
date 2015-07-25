package extract

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/kaffeeshare/data"
)

func imgurl(i *data.Item, sourceURL string, doc *goquery.Document) {
	if !strings.Contains(sourceURL, "imgur.com/") {
		return
	}

	fmt.Println("Running imgurl plugin.")

	selection := doc.Find("meta[property*='og']")

	if selection.Length() != 0 {
		i.ImageURL = ""
		i.Description = ""
		for _, e := range selection.Nodes {
			m := htmlAttributeToMap(e.Attr)

			if m["property"] == "og:image" {
				if !govalidator.IsRequestURL(m["content"]) {
					log.Println("Invalid url in og:image. " + sourceURL)
					continue
				}
				i.Description += "<img src =\""
				i.Description += m["content"]
				i.Description += "\" /><br/>"
			}
		}
	}
}
