// +build !js

package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"

	"github.com/loov/wiki/content/jira"
	"github.com/loov/wiki/content/proxy"
)

func main() {
	http.Handle("/proxy", proxy.New())

	// TODO: use proper test jira instance
	jira, err := jira.New("/jira", "https://issues.apache.org/jira/")
	if err != nil {
		log.Println(err)
	} else {
		http.Handle("/jira/", http.StripPrefix("/jira", jira))
	}

	http.Handle("/data/", http.StripPrefix("/data", http.FileServer(http.Dir("data"))))
	http.Handle("/client/", http.StripPrefix("/client", http.FileServer(http.Dir("client"))))

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/wiki.js", serveFile("wiki.js"))
	http.HandleFunc("/wiki.js.map", serveFile("wiki.js.map"))

	err = http.ListenAndServe("127.0.0.1:8080", nil)
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

	cmd := exec.Command("gopherjs", "build", "github.com/loov/wiki")
	cmd.Env = append(os.Environ(), "GOOS=linux")
	data, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(data))
		log.Println(err)
	}
}
