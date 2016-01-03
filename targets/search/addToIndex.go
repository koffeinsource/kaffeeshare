package search

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/extract"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/search"
)

// DispatchAddToIndex creates the search index for an item passed vie POST
// Used in a Queue.
func DispatchAddToIndex(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if err := r.ParseForm(); err != nil {
		log.Errorf(c, "Error at in DispatchAddToIndex @ ParseForm. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	namespace := r.Form["Namespace"][0]
	URL := r.Form["URL"][0]

	si, err := createSearchItem(c, r.Form, namespace, URL)
	if err != nil {
		log.Errorf(c, "Error while creating data.SearchItem from URL parameters.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = data.AddToSearchIndexTask(c, si, namespace, URL)
	if err != nil {
		log.Errorf(c, "Error while adding the Item to the seach index %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func createSearchItem(c context.Context, values url.Values, namespace string, URL string) (*data.SearchItem, error) {
	DSKey := values["DSKey"][0]

	// We can't do anything without a DSKey
	if DSKey == "" || namespace == "" || URL == "" {
		log.Errorf(c, "Something is wrong with the input data. DSKey: %v, Namespace: %v, URL: %v.", DSKey, namespace, URL)
		return nil, fmt.Errorf("Something is wrong with the input data. DSKey: %v, Namespace: %v, URL: %v.", DSKey, namespace, URL)
	}

	memI, err := memcache.Get(c, DSKey)
	_ = memcache.Delete(c, DSKey) // ignore possible errors

	var searchItem data.SearchItem
	searchItem.Description = values["Caption"][0] + " " + values["Description"][0]

	if err == nil {
		searchItem.HTMLforSearch = search.HTML(string(memI.Value))
	} else {
		// ok, no data in the memcache.
		// we need to re-query the URL to get the HTML data
		i, err := extract.ItemFromURL(URL, c)
		if err != nil {
			log.Errorf(c, "Error in extract.ItemFromURL(). Error: %v", err)
			return nil, err
		}

		searchItem.HTMLforSearch = search.HTML(i.HTMLforSearch)
	}

	if s, err := strconv.ParseInt(values["CreatedAt"][0], 10, 64); err == nil {
		searchItem.CreatedAt = time.Unix(s, 0)
	} else {
		searchItem.CreatedAt = time.Now()
	}

	return &searchItem, nil
}
