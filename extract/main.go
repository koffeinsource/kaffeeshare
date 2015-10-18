package extract

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/extract/plugins"
	"github.com/koffeinsource/kaffeeshare/request"
	"golang.org/x/net/html/charset"
)

// ItemFromURL creates an Item from the passed url
func ItemFromURL(sourceURL string, r *http.Request, log request.Context) (data.Item, error) {

	// Create return value with default values
	returnee := data.Item{
		Caption:   sourceURL,
		URL:       sourceURL,
		CreatedAt: time.Now(),
	}

	// Check if the URL is valid
	// TODO Think if we should handle this differently. Could result in spam?
	if !govalidator.IsRequestURL(sourceURL) {
		errReturn := fmt.Errorf("Invalid URL: %v", sourceURL)
		log.Errorf(errReturn.Error())
		return returnee, errReturn
	}

	contentType, body, err := GetURL(sourceURL, r)
	if err != nil {
		return returnee, err
	}

	//  log.Infof(contentType)
	switch {
	case strings.Contains(contentType, "image/"):
		plugins.Image(&returnee, sourceURL, contentType, log)
	case strings.Contains(contentType, "text/html"):

		var doc *goquery.Document

		charsetReader, err := charset.NewReader(bytes.NewReader(body), contentType)
		if err == nil {
			doc, err = goquery.NewDocumentFromReader(charsetReader)
		} else {
			doc, err = goquery.NewDocumentFromReader(bytes.NewReader(body))
		}

		if err != nil {
			log.Errorf("Problem parsing body. " + sourceURL + " - " + err.Error())
			return returnee, err
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

	return returnee, nil
}
