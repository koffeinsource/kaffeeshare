package share

import (
	"github.com/koffeinsource/go-URLextract"
	"github.com/koffeinsource/kaffeeshare/config"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/httpClient"
	"github.com/koffeinsource/kaffeeshare/search"
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

// CreateURLExtractClient creates a client for go-URLextract
func CreateURLExtractClient(con *data.Context) URLextract.Client {
	var conf URLextract.Client

	conf.HTTPClient = httpClient.Get(con)

	conf.Log = con.Log
	conf.AmazonAdID = config.AmazonAdID
	conf.ImgurClientID = config.ImgurClientID
	return conf
}

// URLsNamespaces shares multiple URLs in mutliple namespaces.
func URLsNamespaces(shareURLs []string, namespaces []string, con *data.Context) error {
	c := CreateURLExtractClient(con)

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

			search.AddToSearchIndex(con, i)
		}
	}

	return errReturn
}
