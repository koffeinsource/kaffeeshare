package extract

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

func getHTTPClient(c context.Context) http.Client {
	s := &urlfetch.Transport{
		Context: c,
		AllowInvalidServerCertificate: true,
	}

	return http.Client{
		Transport: s,
	}
}
