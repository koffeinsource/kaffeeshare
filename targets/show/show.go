package show

import (
	"net/http"
	"text/template"

	"github.com/koffeinsource/notreddit/data"
	"github.com/koffeinsource/notreddit/targets"
	"github.com/koffeinsource/notreddit/targets/startpage"

	"appengine"
)

type wwwTemplateValues struct {
	Namespace string
}

var templateWWW = template.Must(template.ParseFiles("template/base.html", "targets/show/template/html.html"))

//DispatchWWW returns the HTML view of a namespace
func DispatchWWW(w http.ResponseWriter, r *http.Request) {
	//c := appengine.NewContext(r)

	var value wwwTemplateValues
	value.Namespace = targets.GetNamespace(r, "/k/show/www/")

	if value.Namespace == "" {
		startpage.Dispatch(w, r)
		return
	}

	if err := templateWWW.Execute(w, value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(nil)
}

//DispatchRSS returns the rss feed of namespace
// TODO
func DispatchRSS(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// get namespace
	namespace := targets.GetNamespace(r, "/k/show/rss/")
	if namespace == "" {
		startpage.Dispatch(w, r)
		return
	}

	is, _, err := data.GetNewestItems(c, namespace, 20, "")
	if err != nil {
		c.Errorf("Error at in www.dispatch @ GetNewestItem. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Infof("items: %v", is)

	// TODO generate rss
	if err != nil {
		c.Errorf("Error at mashaling in www.dispatch. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(nil)
}
