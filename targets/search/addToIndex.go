package search

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/extract"

	"appengine"
	"appengine/memcache"
	"appengine/search"
)

// DispatchAddToIndex creates the search index for an item passed vie POST
// Used in a Queue.
func DispatchAddToIndex(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if err := r.ParseForm(); err != nil {
		c.Errorf("Error at in DispatchAddToIndex @ ParseForm. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	searchItem := data.ItemSearch{}

	searchItem.DSKey = r.FormValue("DSKey")
	// We can't do anything without a DSKey
	if searchItem.DSKey == "" {
		c.Errorf("There was an error when getting the DSKey")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	namespace := r.FormValue("Namespace")
	if namespace == "" {
		c.Errorf("Got an empty namespace!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	URL := r.FormValue("URL")
	if URL == "" {
		c.Errorf("There was no URL provided")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	index, err := search.Open("items_" + namespace)
	if err != nil {
		c.Errorf("Error while opening the item index %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	memI, err := memcache.Get(c, searchItem.DSKey)
	_ = memcache.Delete(c, searchItem.DSKey) // ignore the error
	if err == nil {
		searchItem.HTMLforSearch = search.HTML(string(memI.Value))
	} else {
		// ok, not data in the memcache.
		// we need to re-query the URL to get the HTML data

		i, err := extract.ItemFromURL(URL, r, c)
		if err != nil {
			c.Errorf("Error in extract.ItemFromURL(). Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		searchItem.HTMLforSearch = search.HTML(i.HTMLforSearch)
	}
	searchItem.HTMLforSearch = search.HTML(optimizeString(string(searchItem.HTMLforSearch)))

	searchItem.Description = r.FormValue("Caption")
	searchItem.Description += " " + r.FormValue("Description")
	searchItem.Description = optimizeString(searchItem.Description)

	if s, err := strconv.ParseInt(r.FormValue("CreatedAt"), 10, 64); err == nil {
		searchItem.CreatedAt = time.Unix(s, 0)
	} else {
		searchItem.CreatedAt = time.Now()
	}

	id, err := index.Put(c, strconv.QuoteToASCII(URL), &searchItem)
	if err != nil {
		c.Errorf("Error while puting the search item in the index %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Debugf("Search item id %v", id)

	w.WriteHeader(http.StatusOK)
}

var htmlScriptRemover = regexp.MustCompile(`<script[^>]*>[\s\S]*?</script>`)
var whiteSpaceCompactor = regexp.MustCompile(`\s+`)

func optimizeString(s string) string {
	s = strings.ToLower(s)
	s = htmlScriptRemover.ReplaceAllString(s, "")
	s = whiteSpaceCompactor.ReplaceAllString(s, " ")
	return s
}
