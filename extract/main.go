package extract

import (
	"bytes"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/koffeinsource/kaffeeshare/data"
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

	contentType, body, err := getURL(sourceURL, r)
	if err != nil {
		log.Errorf(err.Error())
		return returnee
	}

	//  log.Infof(contentType)
	switch {
	case strings.Contains(contentType, "image/"):
		image(&returnee, sourceURL, contentType, log)
	case strings.Contains(contentType, "text/html"):

		// TODO Good check if page is UTF-8 and convert with go-iconv

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Errorf("Problem parsing body. " + sourceURL + " - " + err.Error())
			return returnee
		}

		// Make sure to call this one first
		defaultHTML(&returnee, sourceURL, doc, log)

		amazon(&returnee, sourceURL, doc, log)

		imgurl(&returnee, sourceURL, doc, log)
		gfycat(&returnee, sourceURL, doc, log)

		fefe(&returnee, sourceURL, doc, log)

		youtube(&returnee, sourceURL, doc, log)
		vimeo(&returnee, sourceURL, doc, log)

		dilbert(&returnee, sourceURL, doc, log)
		garfield(&returnee, sourceURL, doc, log)
		xkcd(&returnee, sourceURL, doc, log)
		littlegamers(&returnee, sourceURL, doc, log)

		pastebin(&returnee, sourceURL, doc, log)
	default:
	}

	return returnee
}
