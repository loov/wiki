package mark

import (
	"strings"

	"github.com/loov/wiki/client"
)

type Server struct {
	Host string
}

func NewServer(host string) *Server {
	return &Server{host}
}

func (server *Server) CreateURL(slug string) string {
	return server.Host + slug + ".md"
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
