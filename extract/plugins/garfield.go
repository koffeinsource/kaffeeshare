package plugins

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/kaffeeshare/data"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

// Garfield extracts the comic from a gocomic Garfield page
func Garfield(i *data.Item, sourceURL string, doc *goquery.Document, c context.Context) {
	if !strings.Contains(sourceURL, "www.gocomics.com/garfield") {
		return
	}

	log.Infof(c, "Running Garfield plugin.")

	selection := doc.Find(".strip")
	if len(selection.Nodes) == 0 {
		log.Errorf(c, "Garfield plugin found no .strip. "+sourceURL)
	} else {
		if len(selection.Nodes) > 1 {
			log.Infof(c, "Garfield plugin found >1 .strip. "+sourceURL)
		}
		m := htmlAttributeToMap(selection.Nodes[0].Attr)

		if govalidator.IsRequestURL(m["src"]) {
			i.Description = "<img src =\""
			i.Description += m["src"]
			i.Description += "\" />"
			i.ImageURL = ""
		} else {
			log.Errorf(c, "Garfield plugin invalid url. "+m["src"])
		}
	}

}
