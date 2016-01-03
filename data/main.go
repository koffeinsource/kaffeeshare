package data

import "time"

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

// Data per namespace?
// registered email addresses?
// twitter accounts registered?
