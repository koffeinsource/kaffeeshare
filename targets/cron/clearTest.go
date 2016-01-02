package cron

import (
	"net/http"

	"github.com/koffeinsource/kaffeeshare/data"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
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

	w.WriteHeader(http.StatusOK)
}
