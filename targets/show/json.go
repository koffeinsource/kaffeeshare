package show

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/koffeinsource/kaffeeshare/data"

	"appengine"
)

type jsonReturn struct {
	Items  []data.Item
	Cursor string
}

//DispatchJSON returns the json view of a namespace
func DispatchJSON(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// get namespace
	namespace := mux.Vars(r)["namespace"]
	if namespace == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		c.Errorf("Error at in /show/json @ ParseForm. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cursor := r.FormValue("cursor")

	var returnee jsonReturn
	var err error

	returnee.Items, returnee.Cursor, err = data.GetNewestItems(c, namespace, 20, cursor)
	if err != nil {
		c.Errorf("Error at in /show/json @ GetNewestItem. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Infof("items: %v", returnee.Items)
	c.Infof("cursor: %v", returnee.Cursor)

	s, err := json.Marshal(returnee)
	if err != nil {
		c.Errorf("Error at mashaling in www.dispatch. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(s)
}
