package data

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/search"
	"google.golang.org/appengine/taskqueue"
)

var htmlScriptRemover = regexp.MustCompile(`<script[^>]*>[\s\S]*?</script>`)
var whiteSpaceCompactor = regexp.MustCompile(`\s+`)

func optimizeString(s string) string {
	s = strings.ToLower(s)
	s = htmlScriptRemover.ReplaceAllString(s, "")
	s = whiteSpaceCompactor.ReplaceAllString(s, " ")
	return s
}

// SearchItem is the struct used for the app engine search API
type SearchItem struct {
	DSKey         string
	Description   string
	HTMLforSearch search.HTML
	CreatedAt     time.Time
}

// itemToSearchIndexTask converts a sbset of Item i to url.Values
func (i *Item) itemToSearchIndexTask() url.Values {
	v := url.Values{}

	v.Set("Caption", i.Caption)
	v.Set("Namespace", i.Namespace)
	v.Set("Description", i.Description)
	v.Set("CreatedAt", string(i.CreatedAt.Unix()))
	v.Set("URL", i.URL)
	v.Set("DSKey", i.DSKey)

	return v
}

// AddToSearchIndex adds an Item to the search index.
// To current implementation uses task queues so this operation will
// be executed in the background
func AddToSearchIndex(c context.Context, i Item) {
	// We'll update the search index next
	// FIRST: Store the HTML of the item in the memcache.
	//        We do that because it is often larger than the maximum
	//        task size allowed at the GAE.
	memI := &memcache.Item{
		Key:   i.DSKey,
		Value: []byte(i.HTMLforSearch),
	}
	if err := memcache.Set(c, memI); err != nil {
		log.Infof(c, "Error while storing the search HTML in the memcache for URL %v", i.URL)
	}

	// SECOND: Put the search index update task in the queue
	task := taskqueue.NewPOSTTask("/t/search/add_to_index", i.itemToSearchIndexTask())
	if _, err := taskqueue.Add(c, task, "search-index"); err != nil {
		log.Errorf(c, "Error while triggering the add to index: %v", err)
	} else {
		log.Debugf(c, "Added %v to search-index queue", i.URL)
	}
}

// AddToSearchIndexTask is the implementation of the task described above.
// This adds an item to the search index.
func AddToSearchIndexTask(c context.Context, searchItem *SearchItem, namespace string, URL string) error {
	index, err := search.Open("items_" + namespace)
	if err != nil {
		log.Errorf(c, "Error while opening the item index %v", err)
		return err
	}

	searchItem.HTMLforSearch = search.HTML(optimizeString(string(searchItem.HTMLforSearch)))

	searchItem.Description = optimizeString(searchItem.Description)

	_, err = index.Put(c, strconv.QuoteToASCII(URL), &searchItem)
	if err != nil {
		log.Errorf(c, "Error while puting the search item in the index %v", err)
		return err
	}

	return nil
}
