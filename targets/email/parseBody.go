package email

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"appengine"
)

// TODO update long (>2 characters) top level domains
var urlRegEx = regexp.MustCompile("\\b((http(s?)\\:\\/\\/|~\\/|\\/)|www.)" +
	"(\\w+:\\w+@)?(([-\\w]+\\.)+(com|org|net|gov" +
	"|mil|biz|info|mobi|name|aero|jobs|museum" +
	"|travel|[a-z]{2}))(:[\\d]{1,5})?" +
	"(((\\/([-\\w~!$+|.,=]|%[a-f\\d]{2})+)+|\\/)+|\\?|#)?" +
	"((\\?([-\\w~!$+|.,*:]|%[a-f\\d{2}])+=?" +
	"([-\\w~!$+|.,*:=]|%[a-f\\d]{2})*)" +
	"(&(?:[-\\w~!$+|.,*:]|%[a-f\\d{2}])+=?" +
	"([-\\w~!$+|.,*:=]|%[a-f\\d]{2})*)*)*" +
	"(#([-\\w~!$+|.,*:=]|%[a-f\\d]{2})*)?\\b")

func parseBody(c appengine.Context, mail *email) ([]string, error) {
	if mail.ContentType[:4] == "html" {
		return parseHTMLBody(c, mail.Body)
	}

	if mail.ContentType[:4] == "text" {
		return parseTextBody(c, mail.Body)
	}

	return nil, fmt.Errorf("Unsupported content type: %s", mail.ContentType)
}

func parseHTMLBody(c appengine.Context, body string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	// use a 'set' to remove duplicates
	set := make(map[string]bool)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if !exists {
			return
		}
		set[link] = true

		c.Infof("HTML found %v", link)
	})

	links := make([]string, len(set))
	i := 0
	for k := range set {
		links[i] = k
		i++
	}

	return links, nil
}

func parseTextBody(c appengine.Context, body string) ([]string, error) {
	links := urlRegEx.FindAllString(body, -1)
	c.Infof("Text found %v", links)

	set := make(map[string]bool)
	for _, l := range links {
		set[l] = true
	}

	links = make([]string, len(set))
	i := 0
	for k := range set {
		links[i] = k
		i++
	}

	return links, nil
}
