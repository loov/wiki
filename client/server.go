package client

type Server struct{}

func (server *Server) Load(url string) (*Page, error) {
	return nil, nil
}

func (server *Server) Save(url string, page *Page) error {
	return nil
}
