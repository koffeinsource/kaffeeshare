package extract

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/koffeinsource/kaffeeshare/data"
	"golang.org/x/net/html"
)

// TODO use in most of our code where logging is all we need from a appengine context
type logger interface {
	// Debugf formats its arguments according to the format, analogous to fmt.Printf,
	// and records the text as a log message at Debug level.
	Debugf(format string, args ...interface{})

	// Infof is like Debugf, but at Info level.
	Infof(format string, args ...interface{})

	// Warningf is like Debugf, but at Warning level.
	Warningf(format string, args ...interface{})

	// Errorf is like Debugf, but at Error level.
	Errorf(format string, args ...interface{})

	// Criticalf is like Debugf, but at Critical level.
	Criticalf(format string, args ...interface{})
}

func htmlAttributeToMap(e []html.Attribute) map[string]string {
	m := make(map[string]string)
	for a := range e {
		m[e[a].Key] = e[a].Val
	}
	return m
}

func match(u, startwith string) bool {
	// TODO implement and use in plugins
	return true
}

// ItemFromURL creates an Item from the passed url
func ItemFromURL(sourceURL string, r *http.Request, log logger) data.Item {

	// Create return value with default values
	returnee := data.Item{
		Caption:   sourceURL,
		URL:       sourceURL,
		CreatedAt: time.Now(),
	}

	// TODO refactor the http get code into functions that can be used by the other plugins
	client := getHTTPClient(r)

	// Make a request to the sorceURL
	res, err := client.Get(sourceURL)
	if err != nil {
		log.Infof("Could not get " + sourceURL + " - " + err.Error())
		return returnee
	}
	defer res.Body.Close()

	var body []byte
	// Check the content type of the url
	contentType := res.Header.Get("Content-Type")
	if contentType == "" {
		// No content type send in header. We will try to detect it
		body := make([]byte, 512)
		_, err := res.Body.Read(body)

		if err == nil {
			contentType = http.DetectContentType(body)
		} else {
			// Ok we give up, we cannot access the url
			contentType = "application/octet-stream"
			log.Errorf("Error while reading from the body reader: " + sourceURL + "- " + err.Error())
		}
	}

	//  log.Infof(contentType)
	switch {
	case strings.Contains(contentType, "image/"):
		image(&returnee, sourceURL, contentType, log)
	case strings.Contains(contentType, "text/html"):
		// Read the whole body
		{
			temp, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Errorf("Problem reading the body for " + sourceURL + " - " + err.Error())
				return returnee
			}
			body = append(body, temp...)
		}

		// TODO Good check if page is UTF-8 and convert with go-iconv

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Errorf("Problem parsing body. " + sourceURL + " - " + err.Error())
			return returnee
		}

		// Make sure to call this one first
		defaultHTML(&returnee, sourceURL, doc, log)

		// TODO pass in appengine context for logging!
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
