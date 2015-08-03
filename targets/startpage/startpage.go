package startpage

import (
	"html/template"
	"net/http"

	"github.com/koffeinsource/kaffeeshare/config"
)

var templateWWW = template.Must(template.ParseFiles("template/base.html", "targets/startpage/template/startpage.html"))

type startpageTemplateValues struct {
	URL template.HTML
}

// Dispatch executes all commands for the www target
func Dispatch(w http.ResponseWriter, r *http.Request) {
	value := startpageTemplateValues{URL: config.URL}
	if err := templateWWW.Execute(w, value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
