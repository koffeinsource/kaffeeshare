// +build appengine

package share

import (
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/extract"

	"appengine"
)

// URL shares an URL, i.e. stores it in the datastore and everything
// else that must be done.
func URL(shareURL string, namespace string, c appengine.Context, r *http.Request) error {

	var urls []string
	urls = append(urls, shareURL)

	var namespaces []string
	namespaces = append(namespaces, namespace)

	if err := URLsNamespaces(urls, namespaces, c, r); err != nil {
		return err
	}

	return nil
}

// URLsNamespaces shares multiple URLs in mutliple namespaces.
func URLsNamespaces(shareURLs []string, namespaces []string, c appengine.Context, r *http.Request) error {
	var errReturn error
	errReturn = nil
	for _, shareURL := range shareURLs {
		if !govalidator.IsRequestURL(shareURL) {
			errReturn = fmt.Errorf("Invalid URL: %v", shareURL)
			c.Errorf(errReturn.Error())
			continue
		}

		i := extract.ItemFromURL(shareURL, r, c)

		for _, namespace := range namespaces {
			i.Namespace = namespace
			c.Infof("Sharing item: %v", i)

			if err := data.StoreItem(c, i); err != nil {
				errReturn = err
				c.Errorf("Error at in StoreItem. Item: %v. Error: %v", i, err)
				continue
			}
		}
	}

	return errReturn
}
