package urlextractclient

import (
	"github.com/koffeinsource/go-URLextract"
	"github.com/koffeinsource/kaffeeshare/config"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/httpClient"
)

// Get creates a client for go-URLextract
func Get(con *data.Context) URLextract.Client {
	var conf URLextract.Client

	conf.HTTPClient = httpClient.Get(con)

	conf.Log = con.Log
	conf.AmazonAdID = config.AmazonAdID
	conf.ImgurClientID = config.ImgurClientID
	return conf
}
