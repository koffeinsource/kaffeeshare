package email

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/mail"

	"github.com/koffeinsource/go-imgur"
	"github.com/koffeinsource/kaffeeshare/config"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/share"
)

// used as an return value
type body struct {
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

	con.Log.Debugf("header: %v", msg.Header)

	// get namespaces
	namespaces, err := getNamespaces(msg)
	if err != nil {
		con.Log.Errorf("Error at parsing the receiver fields. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	con.Log.Infof("Detected namespaces: %v", namespaces)

	bodyCopy, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		con.Log.Errorf("Error while reading msg.Body. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get body
	body, err := extractBody(con, msg.Header, bytes.NewBuffer(bodyCopy))
	if err != nil {
		con.Log.Errorf("Error at extracting the body. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	con.Log.Debugf("Received mail: %v", body)

	urls, err := parseBody(con, body)
	if err != nil {
		con.Log.Errorf("Error at parsing the body. Error: %v", err)
	}
	con.Log.Debugf("Found urls: %v", urls)

	if len(urls) > 0 {
		if err := share.URLsNamespaces(urls, namespaces, con); err != nil {
			con.Log.Errorf("Error while sharing URLs. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		images, err := extractAttachment(con, msg.Header, bytes.NewBuffer(bodyCopy))
		if err != nil {
			con.Log.Errorf("Error while extracting attachments. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var imgurclient imgur.Client
		imgurclient.ImgurClientID = config.ImgurClientID
		imgurclient.HTTPClient = share.CreateHTTPClient(con)
		imgurclient.Log = con.Log

		var urls []string

		// FIXME we may set error code more than once!
		for _, im := range images {
			subject := msg.Header.Get("Subject")
			ii, status, err := imgurclient.UploadImage(im.Body, "", im.Encoding, subject, "")
			if status > 399 || err != nil {
				con.Log.Errorf("Error while uploading image. Status: %v Error: %v", status, err)
				w.WriteHeader(http.StatusInternalServerError)
				continue
			}
			var iu data.ImageUpload
			iu.DeleteHash = ii.Deletehash
			iu.URL = ii.Link
			if err := iu.Store(con); err != nil {
				con.Log.Errorf("Error while storing upload in datastore Error: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			urls = append(urls, ii.Link)
		}
		if err := share.URLsNamespaces(urls, namespaces, con); err != nil {
			con.Log.Errorf("Error while sharing URLs. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
