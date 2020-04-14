package main

import (
	"net/http"

	"github.com/koffeinsource/kaffeeshare/targets/check"
	"github.com/koffeinsource/kaffeeshare/targets/cron"
	"github.com/koffeinsource/kaffeeshare/targets/email"
	"github.com/koffeinsource/kaffeeshare/targets/search"
	"github.com/koffeinsource/kaffeeshare/targets/share"
	"github.com/koffeinsource/kaffeeshare/targets/show"
	"github.com/koffeinsource/kaffeeshare/targets/startpage"
	"github.com/koffeinsource/kaffeeshare/targets/update"
	"google.golang.org/appengine"

	"github.com/gorilla/mux"
)

var router = mux.NewRouter()

//<domain>/k/twitter/connect/<namespace>
//<domain>/k/twitter/disconnect/<namespace>
//<domain>/k/email/connect/<namespace>
//<domain>/k/share/<namespace> <- extension url

func init() {
	router.StrictSlash(true)

	router.HandleFunc("/", startpage.Dispatch)
	router.HandleFunc("/k/check/json/{namespace}", check.DispatchJSON)

	// should actually be share/get as we don't do json here
	router.HandleFunc("/k/share/json/{namespace}", share.DispatchJSON)
	router.HandleFunc("/k/share/slack/{namespace}", share.DispatchSlack)

	router.HandleFunc("/k/update/json/{namespace}", update.DispatchJSON)

	router.HandleFunc("/a/search/clear/{namespace}", search.DispatchClearIndex)
	router.HandleFunc("/t/search/clear", search.DispatchClearIndexTask)

	router.HandleFunc("/k/show/json/{namespace}", show.DispatchJSON)
	router.HandleFunc("/k/show/www/{namespace}", show.DispatchWWW)
	router.HandleFunc("/k/show/rss/{namespace}", show.DispatchRSS)

	router.HandleFunc("/c/clear_test/", cron.ClearTest)
	router.HandleFunc("/c/clear_test", cron.ClearTest)

	// TODO move to router
	http.HandleFunc("/_ah/mail/", email.DispatchEmail)
	http.Handle("/", router)
}

func main() {
	/*	http.HandleFunc("/", indexHandler)

		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
			log.Printf("Defaulting to port %s", port)
		}

		log.Printf("Listening on port %s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}*/

	appengine.Main()
}
