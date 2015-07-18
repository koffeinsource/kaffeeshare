package show

import (
	"net/http"

	"github.com/koffeinsource/notreddit/data"
	"github.com/koffeinsource/notreddit/targets"

	"appengine"
)

//DispatchRSS returns the rss feed of namespace
// TODO
func DispatchRSS(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// get namespace
	namespace := targets.GetNamespace(r, "/k/show/rss/")
	if namespace == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	is, _, err := data.GetNewestItems(c, namespace, 20, "")
	if err != nil {
		c.Errorf("Error at in rss.dispatch @ GetNewestItem. Error: %v", err)
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
