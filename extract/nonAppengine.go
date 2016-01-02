// +build !appengine

package extract

/*var internalClient *http.Client

func getHTTPClient(r *http.Request) http.Client {
	var once sync.Once
	body := func() {
		internalClient = &http.Client{}
	}

	once.Do(body)

	return *internalClient
}*/
