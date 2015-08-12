package show

import (
	"net/http"

	"github.com/koffeinsource/kaffeeshare/data"

	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"

	"appengine"
)

//DispatchRSS returns the rss feed of namespace
func DispatchRSS(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	w.Header().Set("Content-Type", "application/rss+xml")

	// get namespace
	namespace := mux.Vars(r)["namespace"]
	if namespace == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cache, err := data.ReadRSSCache(c, namespace)
	if err == nil {
		w.Write([]byte(cache))
		return
	}

	c.Infof("Error at in rss.dispatch while reading the cache. Error: %v", err)

	is, _, err := data.GetNewestItems(c, namespace, 20, "")
	if err != nil {
		c.Errorf("Error at in rss.dispatch @ GetNewestItem. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Infof("items: %v", is)

	feed := &feeds.Feed{
		Title: namespace + " - Kaffeeshare",
		Link:  &feeds.Link{Href: r.URL.String()},
	}

	for _, i := range is {
		rssI := feeds.Item{
			Title:   i.Caption,
			Link:    &feeds.Link{Href: i.URL},
			Created: i.CreatedAt,
		}

		if i.ImageURL != "" {
			rssI.Description += "<div style=\"float:left; margin-right:16px; margin-bottom:16px;\"><img width=\"200\" src=\"" + i.ImageURL + "\" alt=\"\"/></div>"
		}

		rssI.Description += "<p>" + i.Description + "</p><br/><br/>"
		rssI.Description += "<a href=\"" + i.URL + "\">&raquo; " + i.URL + "</a>"

		feed.Items = append(feed.Items, &rssI)
	}

	s, err := feed.ToRss()
	if err != nil {
		c.Errorf("Error at mashaling in www.dispatch. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := data.CacheRSS(c, namespace, s); err != nil {
		c.Errorf("Error at storing the RSS Feed in the cache. Error: %v", err)
	}

	w.Write([]byte(s))
}
