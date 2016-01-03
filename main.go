package main

import (
	"flag"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/kaffeeshare/extract"
	"golang.org/x/net/context"
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
	var c context.Context
	fmt.Print(extract.ItemFromURL(*webaddr, c))
}
