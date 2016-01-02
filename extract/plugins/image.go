package plugins

import (
	"strings"

	"github.com/koffeinsource/kaffeeshare/data"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

// Image is called for links directly to images
func Image(i *data.Item, sourceURL string, contentType string, c context.Context) {
	if !(strings.Index(contentType, "image/") == 0) {
		return
	}

	log.Infof(c, "Running Image plugin.")

	i.ImageURL = ""
	i.Caption = sourceURL[strings.LastIndex(sourceURL, "/")+1:]
	i.Description = "<img src=\"" + sourceURL + "\">"
}
