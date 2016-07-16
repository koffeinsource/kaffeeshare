package cron

import (
	"net/http"

	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/search"
)

// ClearTest deletes all items in the test namespace
func ClearTest(w http.ResponseWriter, r *http.Request) {
	con := data.MakeContext(r)

	// get namespace
	namespace := "test"

	err := data.ClearNamespace(con, namespace)
	if err != nil {
		con.Log.Errorf("Error while clearing the namespace %v: %v", namespace, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	search.ClearSearchItemIndex(con, "test")
	if err != nil {
		con.Log.Errorf("Error while deleting items from search index item_%v: %v", namespace, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
