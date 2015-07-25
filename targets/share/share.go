package share

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/extract"
	"github.com/koffeinsource/kaffeeshare/targets/startpage"

	"appengine"
)

// JSON understood by the extensions
var (
	statusOk    = []byte("{\"status\":\"ok\"}")
	statusError = []byte("{\"status\":\"error\"}")
)

// DispatchJSON receives an extension json request
func DispatchJSON(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// get namespace
	namespace := mux.Vars(r)["namespace"]
	if namespace == "" {
		startpage.Dispatch(w, r)
		return
	}

	shareURL := r.URL.Query().Get("url")

	if !govalidator.IsRequestURL(shareURL) {
		c.Errorf("Error at unmarshalling for share/json. Namespace: %v. Error: %v", namespace, shareURL)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	i := extract.ItemFromURL(shareURL, r)
	i.Namespace = namespace

	c.Infof("Item: %v", i)

	if err := data.StoreItem(c, i); err != nil {
		c.Errorf("Error at in StoreItem. Item: %v. Error: %v", i, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(statusOk)
}
