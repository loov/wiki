package jira

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/gorilla/mux"
)

type Server struct {
	Path    string
	BaseURL string

	client *jira.Client
	router *mux.Router
}

func New(path, baseURL string) (*Server, error) {
	server := &Server{
		Path:    path,
		BaseURL: baseURL,
	}

	var err error
	server.client, err = jira.NewClient(nil, server.BaseURL)
	if err != nil {
		return nil, err
	}

	server.router = mux.NewRouter()
	server.router.HandleFunc("/projects.json", server.Projects)
	server.router.HandleFunc("/p/{projectid}.json", server.Project)
	server.router.HandleFunc("/i/{issueid}.json", server.Issue)

	return server, nil
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.router.ServeHTTP(w, r)
}

func (server *Server) Projects(w http.ResponseWriter, r *http.Request) {
	projects, _, err := server.client.Project.GetList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	page := Page{
		Version:  0,
		Title:    "Projects",
		Modified: time.Now(),
	}

	for _, project := range *projects {
		page.Story = append(page.Story, Item{
			"id":    project.Key,
			"type":  "reference",
			"site":  path.Join("/p", project.Key),
			"title": project.Name,
			"text":  project.Key,
		})
	}

	enc := json.NewEncoder(w)
	if err = enc.Encode(page); err != nil {
		log.Println(err)
	}
}

func (server *Server) Project(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId := vars["projectid"]
	// validate projectId

	project, _, err := server.client.Project.Get(projectId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	page := Page{
		Version:  0,
		Title:    project.Name,
		Modified: time.Now(),
	}

	page.Story = []Item{
		Item{
			"id":   "Description",
			"type": "html",
			"text": project.Description,
		},
	}

	enc := json.NewEncoder(w)
	if err = enc.Encode(page); err != nil {
		log.Println(err)
	}
}

func (server *Server) Issue(w http.ResponseWriter, r *http.Request) {
}
