package email

import (
	"net/http"
	"net/mail"

	"github.com/asaskevich/govalidator"
	"github.com/koffeinsource/notreddit/data"
	"github.com/koffeinsource/notreddit/extract"

	"appengine"
)

// used as an return value
type email struct {
	Body        string
	ContentType string
}

// DispatchEmail parses incoming emails
func DispatchEmail(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	defer r.Body.Close()

	msg, err := mail.ReadMessage(r.Body)
	if err != nil {
		c.Errorf("Error at mail.ReadMessage in DispatchEmail. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Infof("header: %v", msg.Header)

	// get namespaces
	namespaces, err := getNamespaces(msg)
	if err != nil {
		c.Errorf("Error at parsing the receiver fields. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Infof("Detected namespaces: %v", namespaces)

	// get body
	body, err := extractBody(c, msg.Header, msg.Body)
	if err != nil {
		c.Errorf("Error at parsing the body. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Infof("Received mail: %v", body)

	urls, err := parseBody(c, body)
	c.Infof("Found urls: %v", urls)

	for _, shareURL := range urls {
		if !govalidator.IsRequestURL(shareURL) {
			c.Errorf("Invalid URL. Error: %v", shareURL)
			continue
		}

		i := extract.ItemFromURL(shareURL, r)

		for _, namespace := range namespaces {
			i.Namespace = namespace
			c.Infof("Item: %v", i)
		}

		if err := data.StoreItem(c, i); err != nil {
			c.Errorf("Error at in StoreItem. Item: %v. Error: %v", i, err)
			continue
		}
	}

}
