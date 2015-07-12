// +build !appengine

package main

import (
	"flag"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/kaffeeshare2go/extract"
)

func main() {
	webaddr := flag.String("url", "", "The URL to be parsed.")
	flag.Parse()

	// Check if URL was passed at command line
	if *webaddr == "" {
		flag.PrintDefaults()
		return
	}

	// Is it a valid URL?
	if !govalidator.IsRequestURL(*webaddr) {
		fmt.Println("Invalid URL")
		return
	}

	fmt.Print(extract.ItemFromURL(*webaddr, nil))
}
