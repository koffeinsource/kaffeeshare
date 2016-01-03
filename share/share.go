package share

import (
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/extract"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

// URL shares an URL, i.e. stores it in the datastore and everything
// else that must be done.
func URL(shareURL string, namespace string, c context.Context) error {

	var urls []string
	urls = append(urls, shareURL)

	var namespaces []string
	namespaces = append(namespaces, namespace)

	if err := URLsNamespaces(urls, namespaces, c); err != nil {
		return err
	}

	return nil
}

// URLsNamespaces shares multiple URLs in mutliple namespaces.
func URLsNamespaces(shareURLs []string, namespaces []string, c context.Context) error {
	var errReturn error
	errReturn = nil
	for _, shareURL := range shareURLs {

		i, err := extract.ItemFromURL(shareURL, c)
		if err != nil {
			errReturn = err
			log.Errorf(c, "Error in extract.ItemFromURL(). Error: %v", err)
			continue
		}

		for _, namespace := range namespaces {
			i.Namespace = namespace
			//log.Infof(c, "Sharing item: %v", i)

			if err := i.Store(c); err != nil {
				errReturn = err
				log.Errorf(c, "Error in item.Store(). Item: %v. Error: %v", i, err)
				continue
			}

			data.AddToSearchIndex(c, i)
		}
	}

	return errReturn
}
