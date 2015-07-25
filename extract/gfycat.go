package extract

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/koffeinsource/kaffeeshare/data"
	"golang.org/x/net/html"
)

func gfycat(i *data.Item, sourceURL string, doc *goquery.Document) {
	if !strings.Contains(sourceURL, "gfycat.com/") {
		return
	}

	fmt.Println("Running Gfycat plugin.")

	i.ImageURL = ""

	selection := doc.Find(".gfyVid")

	if len(selection.Nodes) == 0 {
		fmt.Println("Gfycat plugin found no .gfyVid. " + sourceURL)
		return
	}
	if len(selection.Nodes) > 1 {
		fmt.Println("Gfycat plugin found >1 .gfyVid. " + sourceURL)
	}
	buf := new(bytes.Buffer)
	err := html.Render(buf, selection.Nodes[0])
	if err != nil {
		fmt.Println("Gfycat plugin error while rendering. " + sourceURL + "- " + err.Error())
		return
	}

	i.Description = buf.String()

	selection = doc.Find(".gfyTitle")
	if len(selection.Nodes) == 0 {
		fmt.Println("Gfycat plugin found no .gfyTitle. " + sourceURL)
		return
	}
	if len(selection.Nodes) > 1 {
		fmt.Println("Gfycat plugin found >1 .gfyTitle. " + sourceURL)
	}
	if len(selection.Nodes) != 0 && selection.Nodes[0].FirstChild != nil {
		i.Caption = selection.Nodes[0].FirstChild.Data
	} else {
		i.Caption = "Gfycat"
	}

}
