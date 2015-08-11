package extract

import (
	"bytes"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/extract/plugins"
	"github.com/koffeinsource/kaffeeshare/request"
)

// ItemFromURL creates an Item from the passed url
func ItemFromURL(sourceURL string, r *http.Request, log request.Context) data.Item {

	// Create return value with default values
	returnee := data.Item{
		Caption:   sourceURL,
		URL:       sourceURL,
		CreatedAt: time.Now(),
	}

	contentType, body, err := GetURL(sourceURL, r)
	if err != nil {
		log.Errorf(err.Error())
		return returnee
	}

	//  log.Infof(contentType)
	switch {
	case strings.Contains(contentType, "image/"):
		plugins.Image(&returnee, sourceURL, contentType, log)
	case strings.Contains(contentType, "text/html"):

		// TODO Good check if page is UTF-8 and convert with go-iconv

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Errorf("Problem parsing body. " + sourceURL + " - " + err.Error())
			return returnee
		}

		// Make sure to call this one first
		plugins.DefaultHTML(&returnee, sourceURL, doc, log)

		plugins.Amazon(&returnee, sourceURL, doc, log)

		plugins.Imgurl(&returnee, sourceURL, doc, log)
		plugins.Gfycat(&returnee, sourceURL, doc, log)

		plugins.Fefe(&returnee, sourceURL, doc, log)

		plugins.Youtube(&returnee, sourceURL, doc, log)
		plugins.Vimeo(&returnee, sourceURL, doc, log)

		plugins.Dilbert(&returnee, sourceURL, doc, log)
		plugins.Garfield(&returnee, sourceURL, doc, log)
		plugins.Xkcd(&returnee, sourceURL, doc, log)
		plugins.Littlegamers(&returnee, sourceURL, doc, log)

		plugins.Pastebin(&returnee, sourceURL, doc, log)
	default:
	}

	return returnee
}
