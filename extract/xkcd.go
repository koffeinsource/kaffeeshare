package extract

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/request"
	"golang.org/x/net/html"
)

func xkcd(i *data.Item, sourceURL string, doc *goquery.Document, log request.Context) {
	if !strings.Contains(sourceURL, "xkcd.com") {
		return
	}

	log.Infof("Running XKCD plugin.")

	selection := doc.Find("#comic")
	if len(selection.Nodes) == 0 {
		log.Infof("XKCD plugin found no #comic. " + sourceURL)
	} else {
		if len(selection.Nodes) > 1 {
			log.Infof("XKCD plugin found >1 #comic. " + sourceURL)
		}
		buf := new(bytes.Buffer)
		err := html.Render(buf, selection.Nodes[0])
		if err != nil {
			log.Infof("XKCD plugin error while rendering. " + sourceURL + "- " + err.Error())
			return
		}
		i.Description = buf.String()
	}

	selection = doc.Find("#ctitle")
	if len(selection.Nodes) == 0 {
		log.Infof("XKCD plugin found no #ctitle. " + sourceURL)
	} else {
		if len(selection.Nodes) > 1 {
			log.Infof("XKCD plugin found >1 #ctitle. " + sourceURL)
		}
		if selection.Nodes[0].FirstChild != nil {
			i.Caption = "XKCD - " + selection.Nodes[0].FirstChild.Data
		} else {
			i.Caption = "XKCD"
		}
	}

	i.ImageURL = ""

}
