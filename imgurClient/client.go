package imgurclient

import (
	"github.com/koffeinsource/go-imgur"
	"github.com/koffeinsource/kaffeeshare/config"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/httpClient"
)

// Get returns an imgurclient
func Get(con *data.Context) *imgur.Client {
	var imgurclient imgur.Client
	imgurclient.ImgurClientID = config.ImgurClientID
	imgurclient.MashapeKey = config.MashapeKey
	imgurclient.HTTPClient = httpClient.Get(con)
	imgurclient.Log = con.Log

	return &imgurclient
}

// GetWithLongDeadline returns an imgurclient with a higher timeout suitable for image upload
func GetWithLongDeadline(con *data.Context) *imgur.Client {
	imgurclient := Get(con)
	imgurclient.HTTPClient = httpClient.GetWithLongDeadline(con)

	return imgurclient
}
