package show

import (
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/koffeinsource/kaffeeshare/config"
	"github.com/koffeinsource/kaffeeshare/targets/startpage"
)

type wwwTemplateValues struct {
	Namespace string
	URL       string
}

var templateWWW = template.Must(template.ParseFiles("template/base.html", "targets/show/template/html.html"))

//DispatchWWW returns the HTML view of a namespace
func DispatchWWW(w http.ResponseWriter, r *http.Request) {
	var value wwwTemplateValues
	value.Namespace = mux.Vars(r)["namespace"]
	value.URL = config.URL

	if value.Namespace == "" {
		startpage.Dispatch(w, r)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=86400") // 1 day
	w.Header().Set("Pragma", "Public")

	if err := templateWWW.Execute(w, value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
