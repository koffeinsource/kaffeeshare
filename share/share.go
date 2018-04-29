package share

import (
	"github.com/koffeinsource/kaffeeshare/URLExtractClient"
	"github.com/koffeinsource/kaffeeshare/data"
)

// URL shares an URL, i.e. stores it in the datastore and everything
// else that must be done.
func URL(shareURL string, namespace string, con *data.Context) error {

	var urls []string
	urls = append(urls, shareURL)

	var namespaces []string
	namespaces = append(namespaces, namespace)

	if err := URLsNamespaces(urls, namespaces, con); err != nil {
		return err
	}

	return nil
}

// URLsNamespaces shares multiple URLs in mutliple namespaces.
func URLsNamespaces(shareURLs []string, namespaces []string, con *data.Context) error {
	if len(shareURLs) == 0 || len(namespaces) == 0 {
		return nil
	}

	c := urlextractclient.Get(con)

	var errReturn error
	errReturn = nil
	for _, shareURL := range shareURLs {
		info, err := c.Extract(shareURL)
		if err != nil {
			errReturn = err
			con.Log.Errorf("Error in URLextract.Extract(). Error: %v", err)
			continue
		}
		var i data.Item
		i = data.ItemFromWebpageInfo(info)

		for _, namespace := range namespaces {
			i.Namespace = namespace
			//con.Log.Infof("Sharing item: %v", i)

			if err := i.Store(con); err != nil {
				errReturn = err
				con.Log.Errorf("Error in item.Store(). Item: %v. Error: %v", i, err)
				continue
			}

		}
	}

	return errReturn
}
