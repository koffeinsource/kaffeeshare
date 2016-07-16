package search

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/search"
)

type returnT struct {
	Items  []data.Item `json:"items"`
	Cursor string      `json:"cursor"`
}

// DispatchSearchJSON searches for a search string in our search index and returns the results as a JSON.
func DispatchSearchJSON(w http.ResponseWriter, r *http.Request) {
	con := data.MakeContext(r)

	namespace := mux.Vars(r)["namespace"]
	if namespace == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		con.Log.Errorf("Error at in k/search/json @ ParseForm. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO cursor support!

	query := r.URL.Query().Get("search")

	con.Log.Infof("Searching for %v in namespace %v", query, namespace)

	var ret returnT
	var err error

	ret.Items, ret.Cursor, err = search.Search(con, namespace, query)
	if err != nil {
		con.Log.Errorf("Error at in k/search/json @ search. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	con.Log.Infof("items: %v", ret.Items)
	con.Log.Infof("cursor: %v", ret.Cursor)

	s, err := json.Marshal(ret)
	if err != nil {
		con.Log.Errorf("Error at mashaling in k/search/json. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(s)
}
