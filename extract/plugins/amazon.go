package plugins

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/kaffeeshare/config"
	"github.com/koffeinsource/kaffeeshare/data"
	"golang.org/x/net/context"
	"golang.org/x/net/html"
	"google.golang.org/appengine/log"
)

// Amazon webpage plugin
func Amazon(i *data.Item, sourceURL string, doc *goquery.Document, c context.Context) {
	if !strings.Contains(sourceURL, "www.amazon.") {
		return
	}

	log.Infof(c, "Running Amazon plugin.")

	// find picture
	{
		selection := doc.Find("#landingImage")
		if len(selection.Nodes) == 0 {
			log.Infof(c, "Amazon plugin found no #landingImage. "+sourceURL)
		} else {
			if len(selection.Nodes) > 1 {
				log.Infof(c, "Amazon plugin found >1 #landingImage. "+sourceURL)
			}
			for _, e := range selection.Nodes {
				if e.Type == html.ElementNode && e.Data == "img" {
					m := htmlAttributeToMap(e.Attr)
					if govalidator.IsRequestURL(m["data-old-hires"]) {
						i.ImageURL = m["data-old-hires"]
					} else {
						log.Infof(c, "Amazon plugin imgURL invalid. "+m["data-old-hires"])
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
			log.Infof(c, "Amazon plugin found no #productTitle. "+sourceURL)
		} else {
			if len(selection.Nodes) > 1 {
				log.Infof(c, "Amazon plugin found >1 #productTitle. "+sourceURL)
			}
			for _, e := range selection.Nodes {
				if e.Type == html.ElementNode && e.Data == "span" {
					i.Caption = e.FirstChild.Data
				}
			}
		}
	}

	// Store HTML for the search
	{
		if s, err := doc.Find(".a-container").Html(); err != nil {
			log.Errorf(c, "Error finding .a-container in HTML: %v", err)
		} else {
			i.HTMLforSearch = s
		}
	}
}
