package fedwiki

import (
	"strings"

	"github.com/loov/wiki/client"
)

type Server struct {
	Host string
}

func NewServer(host string) *Server {
	return &Server{Host: host}
}

func (server *Server) CreateURL(slug string) string {
	return server.Host + "/data/" + slug + ".json"
}

func (server *Server) Open(title, slug string) client.View {
	target := server
	url := slug
	if strings.HasPrefix(slug, "http://") || strings.HasPrefix(slug, "https://") {
		target = server
	} else {
		url = server.CreateURL(slug)
	}
	return NewView(target, title, url)
}
