package search

import (
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/appengine/taskqueue"

	"github.com/gorilla/mux"
	"github.com/koffeinsource/kaffeeshare/data"
	gaesearch "google.golang.org/appengine/search"
)

// DispatchClearIndex triggers an index clear by adding a task to the queue.
func DispatchClearIndex(w http.ResponseWriter, r *http.Request) {
	con := data.MakeContext(r)

	namespace := mux.Vars(r)["namespace"]
	if namespace == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := ClearSearchItemIndex(con, namespace)
	if err != nil {
		con.Log.Errorf("Error at in a/search/clear. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ClearSearchItemIndex removes every entry from an item search index
// Also called directly from the search package after deleting <magic number>
// of docuents.
func ClearSearchItemIndex(con *data.Context, namespace string) error {
	v := url.Values{}

	v.Set("Namespace", namespace)

	task := taskqueue.NewPOSTTask("/t/search/clear", v)
	if _, err := taskqueue.Add(con.C, task, "search-q"); err != nil {
		con.Log.Errorf("Error while adding a clear search index task: %v", err)
		return err
	}

	return nil
}

// DispatchClearIndexTask dispatches the real task of the taskqueue.
// Run in a task queue.
// Enqueued by: ClearSearchItemIndex
func DispatchClearIndexTask(w http.ResponseWriter, r *http.Request) {
	con := data.MakeContext(r)

	if err := r.ParseForm(); err != nil {
		con.Log.Errorf("Error at in DispatchClearIndexTask @ ParseForm. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	namespace := r.Form["Namespace"][0]

	if namespace == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := ClearSearchItemIndexTask(con, namespace)
	if err != nil {
		con.Log.Errorf("Error at in t/search/clear. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ClearSearchItemIndexTask removes every entry from an item search index
func ClearSearchItemIndexTask(con *data.Context, namespace string) error {
	namespace = strings.ToLower(namespace)

	index, err := gaesearch.Open("items_" + namespace)
	if err != nil {
		con.Log.Errorf("Error while opening the item search index %v", err)
		return err
	}

	iter := index.List(con.C, nil)

	counter := 0
	for {

		var id string
		id, err = iter.Next(nil)
		if err == gaesearch.Done {
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
