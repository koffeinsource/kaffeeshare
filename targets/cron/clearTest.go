package cron

import (
	"net/http"

	"github.com/koffeinsource/kaffeeshare/data"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/search"
)

// ClearTest deletes all items in the test namespace
func ClearTest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// get namespace
	namespace := "test"

	err := data.ClearNamespace(c, namespace)
	if err != nil {
		log.Errorf(c, "Error while clearing the namespace %v: %v", namespace, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO clear search index!
	index, err := search.Open("items_" + namespace)
	if err != nil {
		log.Errorf(c, "Error while opening search index item_%v: %v", namespace, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for t := index.List(c, nil); ; {
		id, err := t.Next(nil)
		if err == search.Done {
			break
		}
		if err != nil {
			log.Errorf(c, "Error while deleting items from search index item_%v: %v", namespace, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		index.Delete(c, id)
	}
	w.WriteHeader(http.StatusOK)
}
