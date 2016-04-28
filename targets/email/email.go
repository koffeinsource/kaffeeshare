package email

import (
	"fmt"
	"net/http"

	"github.com/koffeinsource/go-imgur"
	"github.com/koffeinsource/kaffeeshare/config"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/emailparse"
	"github.com/koffeinsource/kaffeeshare/httpClient"
	"github.com/koffeinsource/kaffeeshare/share"
)

// used as an return value
type body struct {
	Body        string
	ContentType string
}

// DispatchEmail handles incoming emails
func DispatchEmail(w http.ResponseWriter, r *http.Request) {
	con := data.MakeContext(r)

	em, err := emailparse.Get(con, r.Body)
	if err != nil {
		con.Log.Errorf("Error at parsing email. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(em.Images) > 0 {
		err2 := handleImages(con, em)
		if err2 != nil {
			con.Log.Errorf("%v", err2)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	if len(em.Texts) > 0 {
		err2 := handleText(con, em)
		if err2 != nil {
			con.Log.Errorf("%v", err2)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	con.Log.Errorf("Found no images and no URLs?.")
	w.WriteHeader(http.StatusInternalServerError)
}

func handleText(con *data.Context, em *emailparse.Email) error {
	urls, err := emailparse.URLsFromText(con, em)
	if err != nil {
		return fmt.Errorf("Error at parsing the body. Error: %v", err)
	}
	con.Log.Debugf("Found urls: %v", urls)

	if err = share.URLsNamespaces(urls, em.Namespaces, con); err != nil {
		return fmt.Errorf("Error while sharing URLs. Error: %v", err)
	}
	return nil
}

func handleImages(con *data.Context, em *emailparse.Email) error {
	var imgurclient imgur.Client
	imgurclient.ImgurClientID = config.ImgurClientID
	imgurclient.HTTPClient = httpClient.GetWithLongDeadline(con)
	imgurclient.Log = con.Log

	var urls []string

	errorHappend := false
	for _, im := range em.Images {
		ii, status, err := imgurclient.UploadImage(im.Body, "", im.Encoding, em.Subject, "")
		if status > 399 || err != nil {
			con.Log.Errorf("Error while uploading image. Status: %v Error: %v", status, err)
			errorHappend = true
			continue
		}
		var iu data.ImageUpload
		iu.DeleteHash = ii.Deletehash
		iu.URL = ii.Link
		if err := iu.Store(con); err != nil {
			con.Log.Errorf("Error while storing upload in datastore Error: %v", err)
			// No important error ...
			// errorHappend = true
		}
		urls = append(urls, ii.Link)
	}
	if err := share.URLsNamespaces(urls, em.Namespaces, con); err != nil {
		return fmt.Errorf("Error while sharing URLs. Error: %v", err)
	}

	if errorHappend {
		return fmt.Errorf("There were some error. We'll retry.")
	}

	return nil
}
