package extract

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/request"
)

func imgurl(i *data.Item, sourceURL string, doc *goquery.Document, log request.Context) {
	if !strings.Contains(sourceURL, "imgur.com/") {
		return
	}

	log.Infof("Running imgurl plugin.")

	selection := doc.Find("meta[property*='og']")

	if selection.Length() != 0 {
		set := make(map[string]bool)

		i.Description = ""
		i.ImageURL = ""

		for _, e := range selection.Nodes {
			m := htmlAttributeToMap(e.Attr)

			if m["property"] == "og:image" {
				if !govalidator.IsRequestURL(m["content"]) {
					log.Infof("Invalid url in og:image. " + sourceURL)
					continue
				}
				if _, in := set[m["content"]]; !in {
					i.Description += "<img src =\""
					temp := strings.Replace(m["content"], "http://", "https://", 1)
					i.Description += temp
					i.Description += "\" /><br/>"
					set[m["content"]] = true
				}
			}
		}

	}
}
