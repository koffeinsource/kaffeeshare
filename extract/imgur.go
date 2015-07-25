package extract

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/notreddit/data"
)

func imgurl(i *data.Item, sourceURL string, doc *goquery.Document) {
	fmt.Println("imgurl! " + sourceURL)
	if !strings.Contains(sourceURL, "imgur.com/") {
		return
	}

	fmt.Println("Running imgurl plugin.")

	// update title
	selection := doc.Find(".image-container")
	if selection.Length() != 0 {
		i.Description += " - ALBUM"
	}

	found := 0
	// get all img within divs with class image
	doc.Find("div.image img").Each(func(unneeded int, s *goquery.Selection) {
		m := htmlAttributeToMap(selection.Nodes[0].Attr)

		if govalidator.IsRequestURL(m["src"]) {
			i.Description = "<img src =\""
			i.Description += m["src"]
			i.Description += "\" /> <br/>"
		} else {
			fmt.Println("imgurl plugin invalid url. " + m["src"])
		}
		i.ImageURL = ""
	})
	if found == 0 {
		fmt.Println("imgurl plugin found no div.image. " + sourceURL)
	}

}
