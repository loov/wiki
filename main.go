package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"sync"
)

func main() {
	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL, _ = url.Parse(r.URL.Query().Get("url"))
		},
	}
	http.Handle("/proxy", proxy)

	http.Handle("/data/", http.StripPrefix("/data", http.FileServer(http.Dir("data"))))
	http.Handle("/client/", http.StripPrefix("/client", http.FileServer(http.Dir("client"))))

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/frontend.js", serveFile("frontend.js"))
	http.HandleFunc("/frontend.js.map", serveFile("frontend.js.map"))

	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	buildClient()
	http.ServeFile(w, r, "index.html")
}

func serveFile(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, name)
	}
}

var mu sync.Mutex

func buildClient() {
	mu.Lock()
	defer mu.Unlock()

	cmd := exec.Command("gopherjs", "build", "github.com/loov/wiki/frontend")
	cmd.Env = append(os.Environ(), "GOOS=linux")
	data, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(data))
		log.Println(err)
	}
}
