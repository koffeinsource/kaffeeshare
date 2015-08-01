package share

import (
	"net/http"
	"text/template"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/extract"
	"github.com/koffeinsource/kaffeeshare/targets/startpage"

	"appengine"
)

var templateFirefox = template.Must(template.ParseFiles("targets/share/template/firefox.html"))

type firefoxTemplateValues struct {
	Message string
}

// DispatchFirefox handels firefox shares
func DispatchFirefox(w http.ResponseWriter, r *http.Request) {
	value := firefoxTemplateValues{Message: "Error! Sorry. Try again"}

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
		if err := templateFirefox.Execute(w, value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	i := extract.ItemFromURL(shareURL, r, c)
	i.Namespace = namespace

	c.Infof("Item: %v", i)

	if err := data.StoreItem(c, i); err != nil {
		c.Errorf("Error at in StoreItem. Item: %v. Error: %v", i, err)
		if err := templateFirefox.Execute(w, value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	value.Message = "Shared"
	if err := templateFirefox.Execute(w, value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
