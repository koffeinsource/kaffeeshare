// +build appengine

package main

import (
	"net/http"

	"github.com/koffeinsource/kaffeeshare2go/targets/check"
	"github.com/koffeinsource/kaffeeshare2go/targets/share"
	"github.com/koffeinsource/kaffeeshare2go/targets/show"
	"github.com/koffeinsource/kaffeeshare2go/targets/startpage"
)

// TODO update. Moved namespace to the end!
//<domain>/k/<namespace>/check/json <- check namespace status
//<domain>/k/<namespace>/show/www <- html ansicht
//<domain>/k/<namespace>/show/rss <- rss feed
//<domain>/k/<namespace>/twitter/connect
//<domain>/k/<namespace>/twitter/disconnect
//<domain>/k/<namespace>/email/connect
//<domain>/k/<namespace>/share <- extension url

func init() {
	http.HandleFunc("/", startpage.Dispatch)
	http.HandleFunc("/k/check/json/", check.DispatchJSON)
	http.HandleFunc("/k/share/json/", share.DispatchJSON)
	http.HandleFunc("/k/show/json/", show.DispatchJSON)
	http.HandleFunc("/k/show/www/", show.DispatchWWW)
	http.HandleFunc("/k/show/rss/", show.DispatchRSS)
}
