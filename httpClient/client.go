package httpClient

import (
	"net/http"
	"time"

	"github.com/koffeinsource/kaffeeshare/data"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

// TODO extract urlfetch.Transport in own function

// Get creates an HTTP client for the GAE
func Get(con *data.Context) *http.Client {
	var timeout time.Time
	timeout = time.Now().Add(60 * time.Second)
	c, _ := context.WithDeadline(con.C, timeout)
	s := &urlfetch.Transport{
		Context: c,
		AllowInvalidServerCertificate: true,
	}
	h := &http.Client{
		Transport: s,
	}
	return h
}

// GetWithLongDeadline creates an HTTP client with a long deadline
func GetWithLongDeadline(con *data.Context) *http.Client {
	var timeout time.Time
	timeout = time.Now().Add(60 * time.Second)
	c, _ := context.WithDeadline(con.C, timeout)
	s := &urlfetch.Transport{
		Context: c,
		AllowInvalidServerCertificate: true,
	}
	h := &http.Client{
		Transport: s,
	}
	return h
}
