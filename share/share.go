package share

import (
	"net/http"

	"github.com/koffeinsource/kaffeeshare/extract"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/taskqueue"
)

// URL shares an URL, i.e. stores it in the datastore and everything
// else that must be done.
func URL(shareURL string, namespace string, c context.Context, r *http.Request) error {

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
func URLsNamespaces(shareURLs []string, namespaces []string, c context.Context, r *http.Request) error {
	var errReturn error
	errReturn = nil
	for _, shareURL := range shareURLs {

		i, err := extract.ItemFromURL(shareURL, r, c)
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

			// We'll update the search index next
			// FIRST: Store the HTML of the item in the memcache.
			//        We do that because it is often larger than the maximum
			//        task size allowed at the GAE.
			{
				memI := &memcache.Item{
					Key:   i.DSKey,
					Value: []byte(i.HTMLforSearch),
				}
				if err := memcache.Set(c, memI); err != nil {
					log.Infof(c, "Error while storing the search HTML in the memcache for URL %v", i.URL)
				}
			}

			// SECOND: Put the search index update task in the queue
			task := taskqueue.NewPOSTTask("/t/search/add_to_index", i.ItemToSearchIndexTask())
			if _, err := taskqueue.Add(c, task, "search-index"); err != nil {
				log.Errorf(c, "Error while triggering the add to index: %v", err)
			} else {
				log.Debugf(c, "Added %v to search-index queue", i.URL)
			}
		}
	}

	return errReturn
}
