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

func (server *Server) CreateURL(href string) string {
	return server.Host + href
}

func (server *Server) Open(title, href string) client.View {
	target := server
	url := href
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		target = server
	} else {
		url = server.CreateURL(href)
	}
	return NewView(target, title, url)
}
