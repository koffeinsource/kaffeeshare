package search

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/search"
	"github.com/koffeinsource/kaffeeshare/share"
	"google.golang.org/appengine/memcache"
	gaesearch "google.golang.org/appengine/search"
)

// DispatchAddToIndex creates the search index for an item passed vie POST
// Used in a Queue.
func DispatchAddToIndex(w http.ResponseWriter, r *http.Request) {
	con := data.MakeContext(r)

	if err := r.ParseForm(); err != nil {
		con.Log.Errorf("Error at in DispatchAddToIndex @ ParseForm. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	namespace := r.Form["Namespace"][0]
	URL := r.Form["URL"][0]

	si, err := createSearchItem(con, r.Form, namespace, URL)
	if err != nil {
		con.Log.Errorf("Error while creating data.Search.Item from URL parameters.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = search.AddToSearchIndexTask(con, si, namespace, URL)
	if err != nil {
		con.Log.Errorf("Error while adding the Item to the seach index %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func createSearchItem(con *data.Context, values url.Values, namespace string, URL string) (*search.Item, error) {
	DSKey := values["DSKey"][0]

	// We can't do anything without a DSKey
	if DSKey == "" || namespace == "" || URL == "" {
		con.Log.Errorf("Something is wrong with the input data. DSKey: %v, Namespace: %v, URL: %v.", DSKey, namespace, URL)
		return nil, fmt.Errorf("Something is wrong with the input data. DSKey: %v, Namespace: %v, URL: %v.", DSKey, namespace, URL)
	}

	memI, err := memcache.Get(con.C, DSKey)
	_ = memcache.Delete(con.C, DSKey) // ignore possible errors

	var searchItem search.Item
	searchItem.Description = values["Caption"][0] + " " + values["Description"][0]

	if err == nil {
		searchItem.HTMLforSearch = gaesearch.HTML(string(memI.Value))
	} else {
		// ok, no data in the memcache.
		// we need to re-query the URL to get the HTML data
		c := share.CreateURLExtractClient(con)
		info, err := c.Extract(URL)
		if err != nil {
			con.Log.Errorf("Error in URLextract.Extract(). Error: %v", err)
			return nil, err
		}
		var i data.Item
		i = data.ItemFromWebpageInfo(info)

		searchItem.HTMLforSearch = gaesearch.HTML(i.HTMLforSearch)
	}

	if s, err := strconv.ParseInt(values["CreatedAt"][0], 10, 64); err == nil {
		searchItem.CreatedAt = time.Unix(s, 0)
	} else {
		searchItem.CreatedAt = time.Now()
	}

	return &searchItem, nil
}
