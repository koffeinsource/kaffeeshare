package data

import "google.golang.org/appengine/datastore"

// An ImageUpload is all the data we store from an uploaded image
type ImageUpload struct {
	URL        string `datastore:"URL,index"`
	DeleteHash string `datastore:"DeleteHash"`
}

// Store stores an ImageUpload in the datastore
func (i *ImageUpload) Store(con *Context) error {
	k := datastore.NewKey(con.C, "ImageUpload", i.URL, 0, nil)
	_, err := datastore.Put(con.C, k, i)
	if err != nil {
		con.Log.Errorf("Error while storing ImageUpload in datastore. ImageUpload: %v. Error: %v", i, err)
		return err
	}
	con.Log.Debugf("Stored ImageUpload %+v", i)

	return nil
}
