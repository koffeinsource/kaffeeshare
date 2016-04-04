package show

import (
	"net/http"
	"time"

	"github.com/koffeinsource/kaffeeshare/config"
	"github.com/koffeinsource/kaffeeshare/data"

	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"

	"google.golang.org/appengine/memcache"
)

//DispatchRSS returns the rss feed of namespace
func DispatchRSS(w http.ResponseWriter, r *http.Request) {
	con := data.MakeContext(r)

	w.Header().Set("Content-Type", "application/rss+xml")
	w.Header().Set("Cache-Control", "public, max-age=1800") // 30 minutes
	w.Header().Set("Pragma", "Public")

	// get namespace
	namespace := mux.Vars(r)["namespace"]
	if namespace == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cache, err := data.ReadRSSCache(con, namespace)
	if err == nil {
		w.Write([]byte(cache))
		return
	}
	if err == memcache.ErrCacheMiss {
		con.Log.Infof("Cache miss for namespace %v", namespace)
	} else {
		con.Log.Errorf("Error at in rss.dispatch while reading the cache. Error: %v", err)
	}

	t := time.Now()
	t = t.Add(-24 * time.Hour * config.RSSTimeRangeinDays)
	is, _, err := data.GetNewestItemsByTime(con, namespace, 100, t, "")
	if err != nil {
		con.Log.Errorf("Error at in rss.dispatch @ GetNewestItem. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	con.Log.Infof("items: %v", is)

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
		con.Log.Errorf("Error at mashaling in www.dispatch. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := data.CacheRSS(con, namespace, s); err != nil {
		con.Log.Errorf("Error at storing the RSS Feed in the cache. Error: %v", err)
	}

	w.Write([]byte(s))
}
