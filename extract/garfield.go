package extract

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/notreddit/data"
)

func garfield(i *data.Item, sourceURL string, doc *goquery.Document) {
	if !strings.Contains(sourceURL, "www.gocomics.com/garfield") {
		return
	}

	fmt.Println("Running Garfield plugin.")

	// update title

	selection := doc.Find(".strip")
	if len(selection.Nodes) == 0 {
		fmt.Println("Garfield plugin found no .strip. " + sourceURL)
	} else {
		if len(selection.Nodes) > 1 {
			fmt.Println("Garfield plugin found >1 .strip. " + sourceURL)
		}
		m := htmlAttributeToMap(selection.Nodes[0].Attr)

		if govalidator.IsRequestURL(m["src"]) {
			i.Description = "<img src =\""
			i.Description += m["src"]
			i.Description += "\" />"
		} else {
			fmt.Println("Amazon plugin invalid url. " + m["src"])
		}
		i.ImageURL = ""
	}

}
