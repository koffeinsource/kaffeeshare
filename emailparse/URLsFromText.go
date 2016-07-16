package emailparse

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/koffeinsource/kaffeeshare/data"

	"github.com/mvdan/xurls"
)

// URLsFromText extracts URLs from the email text bodies
func URLsFromText(con *data.Context, em *Email) ([]string, error) {
	if em == nil {
		return nil, nil
	}

	for _, t := range em.Texts {
		if strings.HasPrefix(t.ContentType, contentTypeHTML) {
			u, err := parseHTMLBody(con, t.Body)
			if err != nil {
				con.Log.Debugf("URLsFromText error while parsing HTML", err)
				continue
			}
			return u, nil
		}

		if strings.HasPrefix(t.ContentType, contentTypeText) {
			u, err := parseTextBody(con, t.Body)
			if err != nil {
				con.Log.Debugf("URLsFromText error while parsing Text", err)
				continue
			}
			return u, nil
		}
	}

	return nil, fmt.Errorf("Could not find an URL in the body.")
}

func parseHTMLBody(con *data.Context, body string) ([]string, error) {
	return firstURLFromHTML(con, body)
}

func firstURLFromHTML(con *data.Context, body string) ([]string, error) {
	if body == "" {
		return nil, nil
	}
	strRdr := strings.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(strRdr)
	if err != nil {
		return nil, err
	}

	var links []string
	found := false

	doc.Find("a").First().Each(func(i int, s *goquery.Selection) {
		if found {
			return
		}
		link, exists := s.Attr("href")
		if !exists {
			return
		}
		if strings.Contains(link, "mailto:") {
			return
		}
		links = append(links, link)
		found = true

		con.Log.Infof("HTML found %v", link)
	})

	return links, nil
}

func parseTextBody(con *data.Context, body string) ([]string, error) {
	return firstURLFromText(con, body)
}

func firstURLFromText(con *data.Context, body string) ([]string, error) {
	var links []string

	// ok, let's only look for URLs with a protocoll explicitly specfied
	// this reduces false positives
	{
		l := xurls.Strict.FindAllString(body, -1)
		for _, s := range l {
			if !strings.Contains(s, "mailto:") {
				links = append(links, s)
				return links, nil
			}
		}
	}

	// found no such URL? Ok, lets try a relaxed query.
	{
		l := xurls.Relaxed.FindAllString(body, -1)
		con.Log.Infof("Found urls in body %v,  %v", body, l)
		for _, s := range l {
			if s != "" && !strings.Contains(s, "mailto:") {
				if !(strings.Contains(s, "http://") || strings.Contains(s, "https://")) {
					s = "http://" + s
				}
				links = append(links, s)
				return links, nil
			}
		}
	}
	return links, nil
}
