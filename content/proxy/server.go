package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func New() http.Handler {
	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL, _ = url.Parse(r.URL.Query().Get("url"))
		},
	}
}
