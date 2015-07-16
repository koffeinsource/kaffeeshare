// +build appengine

package main

import (
	"net/http"

	"github.com/koffeinsource/notreddit/targets/check"
	"github.com/koffeinsource/notreddit/targets/share"
	"github.com/koffeinsource/notreddit/targets/show"
	"github.com/koffeinsource/notreddit/targets/startpage"
)

//<domain>/k/check/json/<namespace> <- check namespace status
//<domain>/k/show/www/<namespace> <- html ansicht
//<domain>/k/show/rss/<namespace> <- rss feed
//<domain>/k/twitter/connect/<namespace>
//<domain>/k/twitter/disconnect/<namespace>
//<domain>/k/email/connect/<namespace>
//<domain>/k/share/<namespace> <- extension url

func init() {
	http.HandleFunc("/", startpage.Dispatch)
	http.HandleFunc("/k/check/json/", check.DispatchJSON)
	http.HandleFunc("/k/share/json/", share.DispatchJSON)
	http.HandleFunc("/k/show/json/", show.DispatchJSON)
	http.HandleFunc("/k/show/www/", show.DispatchWWW)
	http.HandleFunc("/k/show/rss/", show.DispatchRSS)
}
