package startpage

import (
	"html/template"
	"net/http"
)

var templateWWW = template.Must(template.ParseFiles("template/base.html", "targets/startpage/template/startpage.html"))

// Dispatch executes all commands for the www target
func Dispatch(w http.ResponseWriter, r *http.Request) {
	if err := templateWWW.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
