package extract

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// copies every HTML attribute in a map for easier searching
func htmlAttributeToMap(e []html.Attribute) map[string]string {
	m := make(map[string]string)
	for a := range e {
		m[e[a].Key] = e[a].Val
	}
	return m
}

func getURL(sourceURL string, r *http.Request) (string, []byte, error) {
	client := getHTTPClient(r)

	// Make a request to the sorceURL
	res, err := client.Get(sourceURL)
	if err != nil {
		return "", nil, errors.New("Could not get " + sourceURL + " - " + err.Error())
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
			return contentType, nil, nil
		}
	}

	if strings.Contains(contentType, "image/") {
		return contentType, nil, nil
	}

	// Read the whole body
	temp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil, errors.New("Problem reading the body for " + sourceURL + " - " + err.Error())
	}
	body = append(body, temp...)

	return contentType, body, nil

}
