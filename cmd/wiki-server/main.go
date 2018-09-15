package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL, _ = url.Parse(r.URL.Query().Get("url"))
		},
	}
	http.Handle("/proxy", proxy)

	err := http.ListenAndServe("127.0.0.1:8001", nil)
	if err != nil {
		log.Fatal(err)
	}
}
