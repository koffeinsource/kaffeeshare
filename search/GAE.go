package search

import (
	"strconv"

	"github.com/koffeinsource/kaffeeshare/data"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/search"
	"google.golang.org/appengine/taskqueue"
)

// AddToSearchIndex adds an Item to the search index.
// To current implementation uses task queues so this operation will
// be executed in the background
func AddToSearchIndex(con *data.Context, i data.Item) {
	// We'll update the search index next
	// FIRST: Store the HTML of the item in the memcache.
	//        We do that because it is often larger than the maximum
	//        task size allowed at the GAE.
	memI := &memcache.Item{
		Key:   i.DSKey,
		Value: []byte(i.HTMLforSearch),
	}
	if err := memcache.Set(con.C, memI); err != nil {
		con.Log.Infof("Error while storing the search HTML in the memcache for URL %v", i.URL)
	}

	// SECOND: Put the search index update task in the queue
	task := taskqueue.NewPOSTTask("/t/search/add_to_index", itemToSearchIndexTask(i))
	if _, err := taskqueue.Add(con.C, task, "search-index"); err != nil {
		con.Log.Errorf("Error while triggering the add to index: %v", err)
	}
}

// AddToSearchIndexTask is the implementation of the task described above.
// This adds an item to the search index.
func AddToSearchIndexTask(con *data.Context, searchItem *Item, namespace string, URL string) error {
	index, err := search.Open("items_" + namespace)
	if err != nil {
		con.Log.Errorf("Error while opening the item index %v", err)
		return err
	}

	searchItem.HTMLforSearch = search.HTML(optimizeSearchInput(string(searchItem.HTMLforSearch)))

	searchItem.Description = optimizeSearchInput(searchItem.Description)

	_, err = index.Put(con.C, strconv.QuoteToASCII(URL), searchItem)
	if err != nil {
		con.Log.Errorf("Error while puting the search item in the index %v", err)
		return err
	}

	return nil
}
