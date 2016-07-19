package search

import (
	"net/http"
	"net/url"

	"google.golang.org/appengine/taskqueue"

	"github.com/gorilla/mux"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/search"
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
	if _, err := taskqueue.Add(con.C, task, "search-index"); err != nil {
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

	err := search.ClearSearchItemIndexTask(con, namespace)
	if err != nil {
		con.Log.Errorf("Error at in t/search/clear. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
