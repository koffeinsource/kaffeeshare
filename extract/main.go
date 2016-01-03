package extract

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/extract/plugins"

	"golang.org/x/net/context"
	"golang.org/x/net/html/charset"
	"google.golang.org/appengine/log"
)

// ItemFromURL creates an Item from the passed url
func ItemFromURL(sourceURL string, c context.Context) (data.Item, error) {

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
		log.Errorf(c, errReturn.Error())
		return returnee, errReturn
	}

	contentType, body, err := GetURL(sourceURL, c)
	if err != nil {
		return returnee, err
	}

	//  log.Infof(contentType)
	switch {
	case strings.Contains(contentType, "image/"):
		plugins.Image(&returnee, sourceURL, contentType, c)
	case strings.Contains(contentType, "text/html"):

		var doc *goquery.Document

		charsetReader, err := charset.NewReader(bytes.NewReader(body), contentType)
		if err == nil {
			doc, err = goquery.NewDocumentFromReader(charsetReader)
		} else {
			doc, err = goquery.NewDocumentFromReader(bytes.NewReader(body))
		}

		if err != nil {
			log.Errorf(c, "Problem parsing body. "+sourceURL+" - "+err.Error())
			return returnee, err
		}

		// Make sure to call this one first
		plugins.DefaultHTML(&returnee, sourceURL, doc, c)

		plugins.Amazon(&returnee, sourceURL, doc, c)

		plugins.Imgurl(&returnee, sourceURL, doc, c)
		plugins.Gfycat(&returnee, sourceURL, doc, c)

		plugins.Fefe(&returnee, sourceURL, doc, c)

		plugins.Youtube(&returnee, sourceURL, doc, c)
		plugins.Vimeo(&returnee, sourceURL, doc, c)

		plugins.Dilbert(&returnee, sourceURL, doc, c)
		plugins.Garfield(&returnee, sourceURL, doc, c)
		plugins.Xkcd(&returnee, sourceURL, doc, c)
		plugins.Littlegamers(&returnee, sourceURL, doc, c)

		plugins.Pastebin(&returnee, sourceURL, doc, c)
	default:
	}

	return returnee, nil
}
