package show

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/koffeinsource/kaffeeshare/data"
	"google.golang.org/appengine/log"

	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

type jsonReturn struct {
	Items  []data.Item
	Cursor string
}

//DispatchJSON returns the json view of a namespace
func DispatchJSON(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// get namespace
	namespace := mux.Vars(r)["namespace"]
	if namespace == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Errorf(c, "Error at in /show/json @ ParseForm. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cursor := r.FormValue("cursor")

	// no cursos == first elements, may be in the cache
	if cursor == "" {
		cache, err := data.ReadJSONCache(c, namespace)
		if err == nil {
			w.Write(cache)
			return
		}

		if err == memcache.ErrCacheMiss {
			log.Infof(c, "Cache miss for namespace %v", namespace)
		} else {
			log.Errorf(c, "Error at in rss.dispatch while reading the cache. Error: %v", err)
		}
	}

	var returnee jsonReturn
	var err error

	returnee.Items, returnee.Cursor, err = data.GetNewestItems(c, namespace, 20, cursor)
	if err != nil {
		log.Errorf(c, "Error at in /show/json @ GetNewestItem. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof(c, "items: %v", returnee.Items)
	log.Infof(c, "cursor: %v", returnee.Cursor)

	s, err := json.Marshal(returnee)
	if err != nil {
		log.Errorf(c, "Error at mashaling in www.json.dispatch. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// only store values in cache for the first entries
	if cursor == "" {
		if err := data.CacheJSON(c, namespace, s); err != nil {
			log.Errorf(c, "Error at storing the JSON in the cache. Error: %v", err)
		}
	}

	w.Write(s)
}
