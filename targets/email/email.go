package email

import (
	"net/http"
	"net/mail"

	"github.com/koffeinsource/kaffeeshare/share"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// used as an return value
type email struct {
	Body        string
	ContentType string
}

// DispatchEmail parses incoming emails
func DispatchEmail(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	msg, err := mail.ReadMessage(r.Body)
	if err != nil {
		log.Errorf(c, "Error at mail.ReadMessage in DispatchEmail. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof(c, "header: %v", msg.Header)

	// get namespaces
	namespaces, err := getNamespaces(msg)
	if err != nil {
		log.Errorf(c, "Error at parsing the receiver fields. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Infof(c, "Detected namespaces: %v", namespaces)

	// get body
	body, err := extractBody(c, msg.Header, msg.Body)
	if err != nil {
		log.Errorf(c, "Error at parsing the body. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Infof(c, "Received mail: %v", body)

	urls, err := parseBody(c, body)
	log.Infof(c, "Found urls: %v", urls)

	if err := share.URLsNamespaces(urls, namespaces, c, r); err != nil {
		log.Errorf(c, "Error while sharing URLs. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
