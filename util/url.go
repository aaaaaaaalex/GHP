package util

import (
	"net/http"
	"strings"
)

// BaseURL takes the current request and returns the path in a format accepted by the HTML <base> tag
func BaseURL(r *http.Request) string {
	path := r.URL.Path
	if !strings.HasSuffix(path, "/") {
		path = path+ "/"
	}
	return path
}

