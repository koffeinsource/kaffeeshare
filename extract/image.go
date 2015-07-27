package extract

import (
	"strings"

	"github.com/koffeinsource/kaffeeshare/data"
)

func image(i *data.Item, sourceURL string, contentType string, log logger) {
	if !(strings.Index(contentType, "image/") == 0) {
		return
	}

	log.Infof("Running Image plugin.")

	i.ImageURL = ""
	i.Caption = sourceURL[strings.LastIndex(sourceURL, "/")+1:]
	i.Description = "<img src=\"" + sourceURL + "\">"
}
