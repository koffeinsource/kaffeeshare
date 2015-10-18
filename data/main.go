package data

import (
	"net/url"
	"time"

	"appengine/search"
)

// An Item is all the data we store from a website
type Item struct {
	Caption       string    `json:"caption" datastore:"Caption,noindex"`
	URL           string    `json:"url" datastore:"URL,index"`
	Via           string    `json:"via" datastore:"Via,noindex"`
	ImageURL      string    `json:"imageURL" datastore:"ImageURL,noindex"`
	Description   string    `json:"description" datastore:"Description,noindex"`
	CreatedAt     time.Time `json:"createdat" datastore:"CreatedAt,index"`
	Namespace     string    `json:"-" datastore:"Namespace,index"`
	HTMLforSearch string    `json:"-" datastore:"-"`
	DSKey         string    `json:"-" datastore:"-"`
}

// ItemSearch is the struct used for the app engine search API
type ItemSearch struct {
	DSKey         string
	Description   string
	HTMLforSearch search.HTML
	CreatedAt     time.Time
}

// ItemToSearchIndexTask converts a sbset of Item i to url.Values
func (i *Item) ItemToSearchIndexTask() url.Values {
	v := url.Values{}

	v.Set("Caption", i.Caption)
	v.Set("Namespace", i.Namespace)
	v.Set("Description", i.Description)
	v.Set("CreatedAt", string(i.CreatedAt.Unix()))
	v.Set("URL", i.URL)
	v.Set("DSKey", i.DSKey)

	return v
}

// Data per namespace?
// registered email addresses?
// twitter accounts registered?
