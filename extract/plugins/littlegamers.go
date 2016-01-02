package plugins

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/kaffeeshare/data"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

// Littlegamers extract a comic from a littlegamers page
func Littlegamers(i *data.Item, sourceURL string, doc *goquery.Document, c context.Context) {
	if !strings.Contains(sourceURL, "www.little-gamers.com") {
		return
	}

	log.Infof(c, "Running little-gamers plugin.")

	selection := doc.Find("img#comic")

	if len(selection.Nodes) == 0 {
		log.Infof(c, "little-gamers plugin found no img#comic. "+sourceURL)
	} else {
		if len(selection.Nodes) > 1 {
			// that should actually never happen
			log.Errorf(c, "little-gamers plugin found >1 img#comic. ??? "+sourceURL)
		}
		m := htmlAttributeToMap(selection.Nodes[0].Attr)

		if govalidator.IsRequestURL(m["src"]) {
			i.Description = "<img src =\""
			i.Description += m["src"]
			i.Description += "\" />"
			i.ImageURL = ""
		} else {
			log.Errorf(c, "little-gamers plugin invalid url. "+m["src"])
		}
	}

}
