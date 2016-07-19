package search

import (
	"net/url"
	"time"

	"github.com/koffeinsource/kaffeeshare/data"
	"google.golang.org/appengine/search"
)

// Item is the struct used for the app engine search API
// TODO unify with data.Item? see https://cloud.google.com/appengine/docs/go/search/reference#hdr-Fields_and_Facets
type Item struct {
	DSKey         string
	Description   string
	HTMLforSearch search.HTML
	CreatedAt     time.Time
}

// ItemToSearchIndexTask converts a subset of Item i to url.Values
func ItemToSearchIndexTask(i data.Item) url.Values {
	v := url.Values{}

	v.Set("Caption", i.Caption)
	v.Set("Namespace", i.Namespace)
	v.Set("Description", i.Description)
	v.Set("CreatedAt", string(i.CreatedAt.Unix()))
	v.Set("URL", i.URL)
	v.Set("DSKey", i.DSKey)

	return v
}
