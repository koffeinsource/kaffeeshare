package share

import (
	"net/http"

	"github.com/koffeinsource/go-URLextract"
	"github.com/koffeinsource/kaffeeshare/config"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/search"
	"google.golang.org/appengine/urlfetch"
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

// CreateURLExtractConfig creates a config that can be used by go-URLextract
func CreateURLExtractConfig(con *data.Context) URLextract.Config {
	var conf URLextract.Config

	s := &urlfetch.Transport{
		Context: con.C,
		AllowInvalidServerCertificate: true,
	}
	conf.HTTPClient = &http.Client{
		Transport: s,
	}

	conf.Log = con.Log
	conf.AmazonAdID = config.AmazonAdID
	return conf
}

// URLsNamespaces shares multiple URLs in mutliple namespaces.
func URLsNamespaces(shareURLs []string, namespaces []string, con *data.Context) error {

	var errReturn error
	errReturn = nil
	for _, shareURL := range shareURLs {
		info, err := URLextract.Extract(shareURL, CreateURLExtractConfig(con))
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
