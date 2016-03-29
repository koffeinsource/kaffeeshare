package email

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/koffeinsource/kaffeeshare/data"

	"github.com/mvdan/xurls"
)

const (
	contentTypeText  = "text"
	contentTypeHTML  = "html"
	contentTypeMulti = "multipart"
)

func parseBody(con *data.Context, mail *email) ([]string, error) {
	if mail.ContentType[:4] == contentTypeHTML {
		return parseHTMLBody(con, mail.Body)
	}

	if mail.ContentType[:4] == contentTypeText {
		return parseTextBody(con, mail.Body)
	}

	return nil, fmt.Errorf("Unsupported content type: %s", mail.ContentType)
}

func parseHTMLBody(con *data.Context, body string) ([]string, error) {
	return firstURLFromHTML(con, body)
}

func firstURLFromHTML(con *data.Context, body string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	// use a 'set' to remove duplicates
	set := make(map[string]bool)
	doc.Find("a").First().Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if !exists {
			return
		}
		set[link] = true

		con.Log.Infof("HTML found %v", link)
	})

	links := make([]string, len(set))
	i := 0
	for k := range set {
		links[i] = k
		i++
	}

	return links, nil
}

func allURLsFromHTML(con *data.Context, body string) ([]string, error) {
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

		con.Log.Infof("HTML found %v", link)
	})

	links := make([]string, len(set))
	i := 0
	for k := range set {
		links[i] = k
		i++
	}

	return links, nil
}

func parseTextBody(con *data.Context, body string) ([]string, error) {
	return firstURLFromText(con, body)
}

func firstURLFromText(con *data.Context, body string) ([]string, error) {
	links := make([]string, 1)
	links[0] = xurls.Relaxed.FindString(body)
	con.Log.Infof("Found urls in body %v,  %v", body, links)

	return links, nil
}

func allURLsFromText(con *data.Context, body string) ([]string, error) {
	links := xurls.Relaxed.FindAllString(body, -1)
	con.Log.Infof("Found urls in body %v,  %v", body, links)

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
