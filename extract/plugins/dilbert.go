package plugins

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/kaffeeshare/data"

	"golang.org/x/net/context"
	"golang.org/x/net/html"
	"google.golang.org/appengine/log"
)

// Dilbert extracts the comic from a dilbert page
func Dilbert(i *data.Item, sourceURL string, doc *goquery.Document, c context.Context) {
	if !(strings.Contains(sourceURL, "feed.dilbert.com/") ||
		strings.Contains(sourceURL, "dilbert.com/strips/")) {
		return
	}

	log.Infof(c, "Running Dilbert plugin.")

	selection := doc.Find(".STR_Image").Find("img")

	if len(selection.Nodes) == 0 {
		log.Errorf(c, "Dilbert plugin found no .STR_Image. "+sourceURL)
		return
	}

	if len(selection.Nodes) > 1 {
		log.Infof(c, "Dilbert plugin found >1 .STR_Image. "+sourceURL)
	}

	e := selection.Nodes[0]
	if e.Type == html.ElementNode && e.Data == "img" {
		m := htmlAttributeToMap(e.Attr)
		u := ""
		if !strings.Contains(m["src"], "://dilbert.com") {
			u += "https://dilbert.com"
		}

		u += m["src"]
		if govalidator.IsRequestURL(u) {
			i.Description = "<img src=\""
			i.Description += u
			i.Description += "\" />"
		} else {
			log.Errorf(c, "Dilbert plugin invalid url. "+u)
		}

	} else {
		log.Errorf(c, "Dilbert plugin no image tag where we expect one.")
	}

	i.ImageURL = ""
	i.Caption = "Dilbert"
}
