package extract

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kennygrant/sanitize"
	"github.com/koffeinsource/kaffeeshare/data"
	"golang.org/x/net/html"
)

func fefe(i *data.Item, sourceURL string, doc *goquery.Document, log logger) {
	if !strings.Contains(sourceURL, "blog.fefe.de/?ts") {
		return
	}
	log.Infof("Running Fefes Blog plugin.")

	selection := doc.Find("li")

	if len(selection.Nodes) == 0 {
		log.Errorf("Fefes Blog plugin found no li. " + sourceURL)
		return
	}

	if len(selection.Nodes) > 1 {
		log.Infof("Fefes Blog plugin found >1 li. " + sourceURL)
	}

	buf := new(bytes.Buffer)
	err := html.Render(buf, selection.Nodes[0])
	if err != nil {
		log.Errorf("Fefes Blog plugin error while rendering. " + sourceURL + "- " + err.Error())
		return
	}
	i.Description = buf.String()
	start := strings.Index(i.Description, "</a>") + 4
	end := strings.Index(i.Description, "</li>")
	i.Description = i.Description[start:end]

	words := strings.Fields(sanitize.HTML(i.Description))
	i.Caption = ""
	for a := 0; len(i.Caption) < 20 && a < len(words); a++ {
		i.Caption += words[a] + " "
	}
	i.Caption = "Fefes Blog - " + strings.TrimSpace(i.Caption) + "..."
	i.ImageURL = ""
}
