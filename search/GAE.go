package search

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/koffeinsource/kaffeeshare/data"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/search"
	"google.golang.org/appengine/taskqueue"
)

// AddToSearchIndexTask is the implementation of the task described above.
// This adds an item to the search index.
func AddToSearchIndexTask(con *data.Context, searchItem *Item, namespace string, URL string) error {
	namespace = strings.ToLower(namespace)

	index, err := search.Open("items_" + namespace)
	if err != nil {
		con.Log.Errorf("Error while opening the item search index %v", err)
		return err
	}

	searchItem.HTMLforSearch = search.HTML(optimizeSearchInput(string(searchItem.HTMLforSearch)))

	searchItem.Description = optimizeSearchInput(searchItem.Description)

	id := strconv.QuoteToASCII(URL)
	id = strings.Join(strings.Fields(id), "")

	_, err = index.Put(con.C, id, searchItem)
	if err != nil {
		con.Log.Errorf("Error while puting the search item in the index %v", err)
		return err
	}

	return nil
}

// Search for items containing query in the namespace
func Search(con *data.Context, namespace string, query string) ([]data.Item, string, error) {
	query = strings.ToLower(query)
	namespace = strings.ToLower(namespace)

	index, err := search.Open("items_" + namespace)
	if err != nil {
		con.Log.Errorf("Error while opening the item search index %v", err)
		return nil, "", err
	}

	var opt search.SearchOptions
	opt.Limit = 20
	opt.Fields = []string{"DSKey"}
	iter := index.Search(con.C, query, &opt)

	var keys = make([]*datastore.Key, 0, 20)
	counter := 0
	for {
		var i Item
		_, err = iter.Next(&i)
		if err == search.Done {
			break
		}

		con.Log.Debugf("search item: %v", i)

		var k *datastore.Key
		k, err = datastore.DecodeKey(i.DSKey)
		if err != nil {
			con.Log.Errorf("Error decoding key returned from seach index: %v, %v", i.DSKey, err)
			return nil, "", err
		}

		keys = append(keys, k)
		if err != nil {
			con.Log.Errorf("Error fetching next item from search index: %v", err)
			return nil, "", err
		}
		counter++
		if counter == 20 {
			break
		}
	}

	var is = make([]data.Item, 0, 20)

	err = datastore.GetMulti(con.C, keys, is)
	if err != nil {
		con.Log.Errorf("Error fetching items from datastore based on keys got from the search index: %v", err)
		return nil, "", err
	}

	return is, string(iter.Cursor()), nil
}

// ClearSearchItemIndex removes every entry from an item search index
func ClearSearchItemIndex(con *data.Context, namespace string) error {
	v := url.Values{}

	v.Set("Namespace", namespace)

	task := taskqueue.NewPOSTTask("/t/search/clear", v)
	if _, err := taskqueue.Add(con.C, task, "search-index"); err != nil {
		con.Log.Errorf("Error while adding a clear search index task: %v", err)
		return err
	}

	return nil
}

// ClearSearchItemIndexTask removes every entry from an item search index
func ClearSearchItemIndexTask(con *data.Context, namespace string) error {
	namespace = strings.ToLower(namespace)

	index, err := search.Open("items_" + namespace)
	if err != nil {
		con.Log.Errorf("Error while opening the item search index %v", err)
		return err
	}

	iter := index.List(con.C, nil)

	counter := 0
	for {

		var id string
		id, err = iter.Next(nil)
		if err == search.Done {
			break
		}
		if err != nil {
			con.Log.Errorf("Error getting next element from the search index for namespace: %v, %v", namespace, err)
			return err
		}

		err = index.Delete(con.C, id)

		if err != nil {
			con.Log.Errorf("Error while deleting entry from seach index: %v, %v", id, err)
			return err
		}

		counter++
		if counter%20 == 0 {
			con.Log.Debugf("Deleted %v entries in the search index for namespace %v", counter, namespace)
		}
		if counter == 1000 {
			ClearSearchItemIndex(con, namespace)
			break
		}
	}
	con.Log.Infof("Deleted %v entries in the search index for namespace %v", counter, namespace)

	return nil
}
