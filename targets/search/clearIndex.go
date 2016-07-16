package search

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/search"
)

// DispatchClearIndex triggers an index clear
func DispatchClearIndex(w http.ResponseWriter, r *http.Request) {
	con := data.MakeContext(r)

	namespace := mux.Vars(r)["namespace"]
	if namespace == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := search.ClearSearchItemIndex(con, namespace)
	if err != nil {
		con.Log.Errorf("Error at in a/search/clear. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DispatchClearIndexTask dispatches the real task of the taskqueue
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
