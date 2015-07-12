package share

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/kaffeeshare2go/data"
	"github.com/koffeinsource/kaffeeshare2go/extract"
	"github.com/koffeinsource/kaffeeshare2go/targets"
	"github.com/koffeinsource/kaffeeshare2go/targets/startpage"

	"appengine"
)

type shareJSON struct {
	URL string `json:"url"`
}

const (
	statusOk = "ok"
)

type returnJSON struct {
	Status string `json:"status"`
}

// DispatchJSON receives an extension json request
func DispatchJSON(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// get namespace
	namespace := targets.GetNamespace(r, "/k/check/json/")
	if namespace == "" {
		startpage.Dispatch(w, r)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var t shareJSON

	if err := decoder.Decode(&t); err != nil {
		c.Errorf("Error at unmarshalling for share/json. Namespace: %v. Error: %v", namespace, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !govalidator.IsRequestURL(t.URL) {
		c.Errorf("Error at unmarshalling for share/json. Namespace: %v. Error: %v", namespace, t.URL)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	i := extract.ItemFromURL(t.URL, r)
	i.Namespace = namespace

	c.Infof("Item: %v", i)

	if err := data.StoreItem(c, i); err != nil {
		c.Errorf("Error at in StoreItem. Item: %v. Error: %v", i, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var returnee returnJSON
	returnee.Status = statusOk

	s, _ := json.Marshal(returnee)

	w.Write(s)
}
