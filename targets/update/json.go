package update

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/koffeinsource/kaffeeshare/data"

	"appengine"
)

type jsonReturn struct {
	LastUpdate int64 `json:"last_update"`
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

	is, _, err := data.GetNewestItems(c, namespace, 1, "")
	if err != nil {
		c.Errorf("Error while getting 1 item for update/json. Error %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var returnee jsonReturn
	if len(is) == 0 {
		returnee.LastUpdate = time.Date(1982, time.May, 14, 0, 0, 0, 0, time.UTC).Unix()
	} else {
		returnee.LastUpdate = is[0].CreatedAt.Unix()
	}

	s, err := json.Marshal(returnee)
	if err != nil {
		c.Errorf("Error at mashaling in update/json. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(s)
}
