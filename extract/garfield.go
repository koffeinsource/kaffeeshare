package extract

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/kaffeeshare/data"
)

func garfield(i *data.Item, sourceURL string, doc *goquery.Document, log logger) {
	if !strings.Contains(sourceURL, "www.gocomics.com/garfield") {
		return
	}

	log.Infof("Running Garfield plugin.")

	selection := doc.Find(".strip")
	if len(selection.Nodes) == 0 {
		log.Errorf("Garfield plugin found no .strip. " + sourceURL)
	} else {
		if len(selection.Nodes) > 1 {
			log.Infof("Garfield plugin found >1 .strip. " + sourceURL)
		}
		m := htmlAttributeToMap(selection.Nodes[0].Attr)

		if govalidator.IsRequestURL(m["src"]) {
			i.Description = "<img src =\""
			i.Description += m["src"]
			i.Description += "\" />"
			i.ImageURL = ""
		} else {
			log.Errorf("Garfield plugin invalid url. " + m["src"])
		}
	}

}
