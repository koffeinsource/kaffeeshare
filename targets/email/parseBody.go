package email

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/koffeinsource/kaffeeshare/request"
	"github.com/mvdan/xurls"
)

func parseBody(c request.Context, mail *email) ([]string, error) {
	if mail.ContentType[:4] == "html" {
		return parseHTMLBody(c, mail.Body)
	}

	if mail.ContentType[:4] == "text" {
		return parseTextBody(c, mail.Body)
	}

	return nil, fmt.Errorf("Unsupported content type: %s", mail.ContentType)
}

func parseHTMLBody(c request.Context, body string) ([]string, error) {
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

func parseTextBody(c request.Context, body string) ([]string, error) {
	links := xurls.Relaxed.FindAllString(body, -1)
	c.Infof("Found urls in body %v,  %v", body, links)

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
