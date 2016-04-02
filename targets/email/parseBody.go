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

func parseBody(con *data.Context, mail *body) ([]string, error) {
	if mail.ContentType[:len(contentTypeHTML)] == contentTypeHTML {
		return parseHTMLBody(con, mail.Body)
	}

	if mail.ContentType[:len(contentTypeText)] == contentTypeText {
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
	l := xurls.Relaxed.FindAllString(body, -1)
	con.Log.Infof("Found urls in body %v,  %v", body, l)
	for _, s := range l {
		if s != "" && !strings.Contains(s, "mailto:") {
			links = append(links, s)
			return links, nil
		}
	}

	return links, nil
}
