// +build appengine

package extract

import (
	"net/http"
	"sync"

	"appengine"
	"appengine/urlfetch"
)

var internalClient *http.Client

func getHTTPClient(r *http.Request) http.Client {
	var once sync.Once

	once.Do(func() {
		c := appengine.NewContext(r)
		s := &urlfetch.Transport{
			Context: c,
			AllowInvalidServerCertificate: true,
		}
		internalClient = &http.Client{
			Transport: s,
		}
	})

	return *internalClient
}
