// +build !appengine

package main

import (
	"flag"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/notreddit/extract"
)

func main() {
	webaddr := flag.String("url", "", "The URL to be parsed.")
	flag.Parse()

	// Check if URL was passed at command line
	if *webaddr == "" {
		flag.PrintDefaults()
		return
	}

	fmt.Print("Validating URL...")
	// Is it a valid URL?
	if !govalidator.IsRequestURL(*webaddr) {
		fmt.Println("invalid!")
		return
	}
	fmt.Println(" success!")

	fmt.Print(extract.ItemFromURL(*webaddr, nil))
}
