package targets

import (
	"net/http"
	"strings"
)

// GetNamespace returns the namespace provided in a http.Request. "" is returned
// in case there is not namespace provided.
func GetNamespace(r *http.Request, baseURL string) string {
	// remove trailing /
	url := strings.ToLower(r.URL.String())
	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}

	// get namespace
	indexNamespace := strings.Index(url, baseURL)
	if indexNamespace == -1 {
		return ""
	}
	return url[indexNamespace+len(baseURL)+1:]
}
