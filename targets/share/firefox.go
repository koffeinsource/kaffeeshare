package share

import (
	"net/http"
	"text/template"
)

var templateFirefox = template.Must(template.ParseFiles("targets/share/template/firefox.html"))

// DispatchFirefox handels firefox shares
func DispatchFirefox(w http.ResponseWriter, r *http.Request) {
	if err := templateFirefox.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
