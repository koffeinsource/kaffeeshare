package email

import (
	"net/http"
	"net/mail"

	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/share"
)

// used as an return value
type email struct {
	Body        string
	ContentType string
}

// DispatchEmail parses incoming emails
func DispatchEmail(w http.ResponseWriter, r *http.Request) {
	con := data.MakeContext(r)

	msg, err := mail.ReadMessage(r.Body)
	if err != nil {
		con.Log.Errorf("Error at mail.ReadMessage in DispatchEmail. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	con.Log.Infof("header: %v", msg.Header)

	// get namespaces
	namespaces, err := getNamespaces(msg)
	if err != nil {
		con.Log.Errorf("Error at parsing the receiver fields. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	con.Log.Infof("Detected namespaces: %v", namespaces)

	// get body
	body, err := extractBody(con, msg.Header, msg.Body)
	if err != nil {
		con.Log.Errorf("Error at extracting the body. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	con.Log.Infof("Received mail: %v", body)

	urls, err := parseBody(con, body)
	if err != nil {
		con.Log.Errorf("Error at parsing the body. Error: %v", err)
	}
	con.Log.Infof("Found urls: %v", urls)

	if err := share.URLsNamespaces(urls, namespaces, con); err != nil {
		con.Log.Errorf("Error while sharing URLs. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
