package extract

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/koffeinsource/notreddit/data"
	"golang.org/x/net/html"
)

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
func ItemFromURL(sourceURL string, r *http.Request) data.Item {

	// Create return value with default values
	returnee := data.Item{
		Caption:   sourceURL,
		URL:       sourceURL,
		CreatedAt: time.Now(),
	}

	client := getHTTPClient(r)

	// Make a request to the sorceURL
	res, err := client.Get(sourceURL)
	if err != nil {
		log.Println("Could not get " + sourceURL + " - " + err.Error())
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
			log.Println("Error while reading from the body reader: " + sourceURL + "- " + err.Error())
		}
	}

	// Read the whole body
	{
		temp, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println("Problem reading the body for " + sourceURL + " - " + err.Error())
			return returnee
		}
		body = append(body, temp...)
	}

	//log.Println(contentType)
	switch {
	case strings.Contains(contentType, "image/"):
		image(&returnee, sourceURL, contentType)
	case strings.Contains(contentType, "text/html"):
		// TODO Good check if page is UTF-8 and convert with go-iconv

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Println("Problem parsing body. " + sourceURL + " - " + err.Error())
			return returnee
		}

		// Make sure to call this one first
		defaultHTML(&returnee, sourceURL, doc)

		// TODO pass in appengine context for logging!
		amazon(&returnee, sourceURL, doc)
		imgurl(&returnee, sourceURL, doc)
		dilbert(&returnee, sourceURL, doc)
		fefe(&returnee, sourceURL, doc)
		garfield(&returnee, sourceURL, doc)
		gfycat(&returnee, sourceURL, doc)
		xkcd(&returnee, sourceURL, doc)
		youtube(&returnee, sourceURL, doc)
		vimeo(&returnee, sourceURL, doc)

	default:
	}

	return returnee
}
