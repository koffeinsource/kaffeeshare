package show

import (
	"encoding/json"
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

//DispatchJSON returns the json view of a namespace
//TODO
func DispatchJSON(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// get namespace
	namespace := targets.GetNamespace(r, "/k/show/json/")
	if namespace == "" {
		startpage.Dispatch(w, r)
		return
	}

	i, err := data.GetNewestItems(c, namespace, 20)
	if err != nil {
		c.Errorf("Error at in www.dispatch @ GetNewestItem. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Infof("items: %v", i)

	s, err := json.Marshal(i)
	if err != nil {
		c.Errorf("Error at mashaling in www.dispatch. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(s)
}

//DispatchRSS returns the rss feed of namespace
//TODO
func DispatchRSS(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// get namespace
	namespace := targets.GetNamespace(r, "/k/show/rss/")
	if namespace == "" {
		startpage.Dispatch(w, r)
		return
	}

	i, err := data.GetNewestItems(c, namespace, 20)
	if err != nil {
		c.Errorf("Error at in www.dispatch @ GetNewestItem. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Infof("items: %v", i)

	// TODO generate rss
	if err != nil {
		c.Errorf("Error at mashaling in www.dispatch. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(nil)
}
