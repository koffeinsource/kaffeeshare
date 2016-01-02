package plugins

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/kaffeeshare/data"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

// Imgurl extract all images from an imgurl album
func Imgurl(i *data.Item, sourceURL string, doc *goquery.Document, c context.Context) {
	if !strings.Contains(sourceURL, "imgur.com/") {
		return
	}

	log.Infof(c, "Running imgurl plugin.")

	selection := doc.Find("meta[property*='og']")

	if selection.Length() != 0 {
		set := make(map[string]bool)

		i.Description = ""
		i.ImageURL = ""

		for _, e := range selection.Nodes {
			m := htmlAttributeToMap(e.Attr)

			if m["property"] == "og:image" {
				if !govalidator.IsRequestURL(m["content"]) {
					log.Infof(c, "Invalid url in og:image. "+sourceURL)
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
