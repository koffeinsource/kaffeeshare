package plugins

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/kaffeeshare/config"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/request"
	"golang.org/x/net/html"
)

// Amazon webpage plugin
func Amazon(i *data.Item, sourceURL string, doc *goquery.Document, log request.Context) {
	if !strings.Contains(sourceURL, "www.amazon.") {
		return
	}

	log.Infof("Running Amazon plugin.")

	// find picture
	{
		selection := doc.Find("#landingImage")
		if len(selection.Nodes) == 0 {
			log.Infof("Amazon plugin found no #landingImage. " + sourceURL)
		} else {
			if len(selection.Nodes) > 1 {
				log.Infof("Amazon plugin found >1 #landingImage. " + sourceURL)
			}
			for _, e := range selection.Nodes {
				if e.Type == html.ElementNode && e.Data == "img" {
					m := htmlAttributeToMap(e.Attr)
					if govalidator.IsRequestURL(m["data-old-hires"]) {
						i.ImageURL = m["data-old-hires"]
					} else {
						log.Infof("Amazon plugin imgURL invalid. " + m["data-old-hires"])
					}
				}
			}
		}
	}

	// update url to contain tag
	{
		// This is our tag. We should make it configurable
		urlExtension := "tag=" + config.AmazonAdID
		start := strings.Index(i.URL, "tag=")
		if start != -1 {
			end := strings.Index(i.URL[start+1:], "&") + start + 1
			i.URL = i.URL[:start] + i.URL[end:]
		}

		if strings.Index(i.URL, "?") == -1 {
			i.URL += "?" + urlExtension
		} else {
			i.URL += "&" + urlExtension
		}
	}

	// update title
	{
		selection := doc.Find("#productTitle")
		if len(selection.Nodes) == 0 {
			log.Infof("Amazon plugin found no #productTitle. " + sourceURL)
		} else {
			if len(selection.Nodes) > 1 {
				log.Infof("Amazon plugin found >1 #productTitle. " + sourceURL)
			}
			for _, e := range selection.Nodes {
				if e.Type == html.ElementNode && e.Data == "span" {
					i.Caption = e.FirstChild.Data
				}
			}
		}

	}
}
