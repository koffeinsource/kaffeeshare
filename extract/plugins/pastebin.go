package plugins

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/koffeinsource/kaffeeshare/data"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

// Pastebin extracts the content from a pastbin page
func Pastebin(i *data.Item, sourceURL string, doc *goquery.Document, c context.Context) {
	if !strings.Contains(sourceURL, "pastebin.com/") {
		return
	}

	// TODO replace the logic below with a query to http://pastebin.com/raw.php?i=
	log.Infof(c, "Running pastebin plugin.")

	selection := doc.Find("#paste_code")

	if len(selection.Nodes) == 0 {
		log.Infof(c, "Pastebin plugin found no #paste_code. "+sourceURL)
	} else {
		if len(selection.Nodes) > 1 {
			log.Infof(c, "Pastebin plugin found >1 #paste_code. "+sourceURL)
		}

		str, err := selection.Html()
		if err != nil {
			log.Infof(c, "Error when creating html in pastebin plugin: ")
			return
		}
		str = strings.Replace(str, "\n", "<br />\n", -1)
		i.Description = str
		i.ImageURL = ""
	}

}
