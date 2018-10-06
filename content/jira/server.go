package jira

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/kr/pretty"

	jira "github.com/andygrunwald/go-jira"
	"github.com/gorilla/mux"
)

type Server struct {
	Path    string
	BaseURL string

	client *jira.Client
	router *mux.Router
}

func NewPublicReadOnly(path, baseURL string) (*Server, error) {
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

func NewBasicAuth(path, baseURL string, user, pass string) (*Server, error) {
	server := &Server{
		Path:    path,
		BaseURL: baseURL,
	}

	tp := jira.BasicAuthTransport{
		Username: user,
		Password: pass,
	}

	var err error
	server.client, err = jira.NewClient(tp.Client(), server.BaseURL)
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
			"id":        project.Key,
			"type":      "reference",
			"site":      path.Join("/p", project.Key),
			"title":     project.Name,
			"text":      project.Key,
			"thumbnail": project.AvatarUrls.Four8X48,
		})
	}

	enc := json.NewEncoder(w)
	if err = enc.Encode(page); err != nil {
		log.Println(err)
	}
}

func (server *Server) Project(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	procjectid := vars["projectid"]
	//TODO: validate procjectid

	project, _, err := server.client.Project.Get(procjectid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	issues, _, err := server.client.Issue.Search("", nil)
	if err != nil {
		page.Story = append(page.Story, Item{
			"id":   "Error",
			"type": "paragraph",
			"text": err.Error(),
		})
	} else {
		for _, issue := range issues {
			page.Story = append(page.Story, Item{
				"id":        issue.Key,
				"type":      "reference",
				"site":      path.Join("/i", issue.Key),
				"title":     issue.Fields.Summary,
				"text":      issue.Key,
				"thumbnail": issue.Fields.Type.IconURL,
			})
		}
	}

	enc := json.NewEncoder(w)
	if err = enc.Encode(page); err != nil {
		log.Println(err)
	}
}

func (server *Server) Issue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	issueid := vars["issueid"]
	//TODO: validate issueid

	issue, _, err := server.client.Issue.Get(issueid, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	page := Page{
		Version:  0,
		Title:    issue.Fields.Summary,
		Class:    "wide",
		Modified: time.Now(),
	}

	fieldsjson, _ := issue.Fields.MarshalJSON()
	renderedfieldsjson := pretty.Sprint(issue.RenderedFields)

	page.Add(Item{
		"id":   "Description",
		"type": "paragraph",
		"text": issue.Fields.Description,
	})

	if issue.Fields.Comments != nil && len(issue.Fields.Comments.Comments) > 0 {
		page.Add(Item{"id": "Activity", "type": "html", "text": "<h2>Activity</h2>"})
		for _, comment := range issue.Fields.Comments.Comments {
			page.Add(Item{
				"id":   "Attachment-" + comment.ID,
				"type": "paragraph",
				"text": comment.Body,
			})
		}
	}

	if len(issue.Fields.Attachments) > 0 {
		page.Add(Item{"id": "Attachments", "type": "html", "text": "<h2>Attachments</h2>"})
		for _, attachment := range issue.Fields.Attachments {
			page.Add(Item{
				"id":        "Attachment-" + attachment.ID,
				"type":      "reference",
				"site":      attachment.Content,
				"title":     attachment.Filename,
				"text":      attachment.Created,
				"thumbnail": attachment.Thumbnail,
			})
		}
	}

	page.Add(Item{
		"id":   "Fields",
		"type": "html",
		"text": `<pre style="white-space:pre-wrap">` + string(fieldsjson) + `</pre>`,
	})

	page.Add(Item{
		"id":   "RenderedFields",
		"type": "html",
		"text": `<pre style="white-space:pre-wrap">` + string(renderedfieldsjson) + `</pre>`,
	})

	enc := json.NewEncoder(w)
	if err = enc.Encode(page); err != nil {
		log.Println(err)
	}
}
