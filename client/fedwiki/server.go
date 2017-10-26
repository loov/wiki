package fedwiki

import "github.com/loov/wiki/client"

type Server struct{}

func (server *Server) Open(title, url string) client.View {
	return NewView(title, url)
}
