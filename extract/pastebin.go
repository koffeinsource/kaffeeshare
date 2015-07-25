package extract

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/koffeinsource/kaffeeshare/data"
)

func pastebin(i *data.Item, sourceURL string, doc *goquery.Document) {
	if !strings.Contains(sourceURL, "pastebin.com/") {
		return
	}

	fmt.Println("Running pastebin plugin.")

	selection := doc.Find("#paste_code")

	if len(selection.Nodes) == 0 {
		fmt.Println("Pastebin plugin found no #paste_code. " + sourceURL)
	} else {
		if len(selection.Nodes) > 1 {
			fmt.Println("Pastebin plugin found >1 #paste_code. " + sourceURL)
		}

		str, err := selection.Html()
		if err != nil {
			fmt.Println("Error when creating html in pastebin plugin: ")
			return
		}
		str = strings.Replace(str, "\n", "<br />\n", -1)
		i.Description = str
		i.ImageURL = ""
	}

}
