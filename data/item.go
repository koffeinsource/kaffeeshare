package data

import (
	"time"

	"github.com/koffeinsource/go-URLextract/webpage"
)

// An Item is all the data we store from a website
// TODO check if URL must be indexed .Namespace?
type Item struct {
	Caption       string    `json:"caption" datastore:"Caption,noindex"`
	URL           string    `json:"url" datastore:"URL,index"`
	Via           string    `json:"via" datastore:"Via,noindex"` // Required for DS backward compatibility
	ImageURL      string    `json:"imageURL" datastore:"ImageURL,noindex"`
	Description   string    `json:"description" datastore:"Description,noindex"`
	CreatedAt     time.Time `json:"createdat" datastore:"CreatedAt,index"`
	Namespace     string    `json:"-" datastore:"Namespace,index"`
	HTMLforSearch string    `json:"-" datastore:"-"` // Required for DS backward compatibility
	DSKey         string    `json:"-" datastore:"-"`
}

// ItemFromWebpageInfo converts the results returned from go-URLextract into our own Item struct
func ItemFromWebpageInfo(info webpage.Info) Item {
	var ret Item
	ret.Caption = info.Caption
	ret.Description = info.Description
	ret.ImageURL = info.ImageURL
	ret.URL = info.URL
	ret.CreatedAt = time.Now()
	return ret
}

// Data per namespace?
// registered email addresses?
// twitter accounts registered?
