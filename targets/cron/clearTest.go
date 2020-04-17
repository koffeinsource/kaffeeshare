package cron

import (
	"net/http"

	"github.com/koffeinsource/kaffeeshare/data"
)

// ClearTest deletes all items in the test namespace
func ClearTest(w http.ResponseWriter, r *http.Request) {
	con := data.MakeContext(r)

	// get namespace
	namespaces := []string{"test", "temp"}

	for _, namespace := range namespaces {
		err := data.ClearNamespace(con, namespace)
		if err != nil {
			con.Log.Errorf("Error while clearing the namespace %v: %v", namespace, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
