package extract

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/koffeinsource/kaffeeshare2go/data"
)

func vimeo(i *data.Item, sourceURL string, doc *goquery.Document) {
	if !strings.Contains(sourceURL, "vimeo.com") {
		return
	}

	fmt.Println("Running Vimeo plugin.")

	// remove trailing '/' of the url, if any
	if string(sourceURL[len(sourceURL)-1]) == "/" {
		sourceURL = sourceURL[:len(sourceURL)-1]
	}
	videoIDstart := strings.LastIndex(sourceURL, "/")
	if videoIDstart == -1 {
		fmt.Println("Vimeo plugin found no '/' ??? " + sourceURL)
		return
	}

	videoIDstart++
	videoID := sourceURL[videoIDstart:]
	i.Description += "<br/><br/><br/><iframe src=\"http://player.vimeo.com/video/"
	i.Description += videoID
	i.Description += "?title=0&amp;byline=0&amp;portrait=0\" width=\"400\" height=\"225\" frameborder=\"0\" webkitAllowFullScreen mozallowfullscreen allowFullScreen></iframe>"

	i.ImageURL = ""
}
