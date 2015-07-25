package extract

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/kaffeeshare/data"
	"golang.org/x/net/html"
)

func amazon(i *data.Item, sourceURL string, doc *goquery.Document) {
	if !strings.Contains(sourceURL, "www.amazon.") {
		return
	}

	fmt.Println("Running Amazon plugin.")

	// find picture
	{
		selection := doc.Find("#landingImage")
		if len(selection.Nodes) == 0 {
			fmt.Println("Amazon plugin found no #landingImage. " + sourceURL)
		} else {
			if len(selection.Nodes) > 1 {
				fmt.Println("Amazon plugin found >1 #landingImage. " + sourceURL)
			}
			for _, e := range selection.Nodes {
				if e.Type == html.ElementNode && e.Data == "img" {
					m := htmlAttributeToMap(e.Attr)
					if govalidator.IsRequestURL(m["data-old-hires"]) {
						i.ImageURL = m["data-old-hires"]
					} else {
						fmt.Println("Amazon plugin imgURL invalid. " + m["data-old-hires"])
					}
				}
			}
		}
	}

	// update url to contain tag
	{
		// This is our tag. We should make it configurable
		urlExtension := "tag=" + "gschaftshuonl-21"
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
			fmt.Println("Amazon plugin found no #productTitle. " + sourceURL)
		} else {
			if len(selection.Nodes) > 1 {
				fmt.Println("Amazon plugin found >1 #productTitle. " + sourceURL)
			}
			for _, e := range selection.Nodes {
				if e.Type == html.ElementNode && e.Data == "span" {
					i.Caption = e.FirstChild.Data
				}
			}
		}

	}
}
