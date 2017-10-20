package client

type Page struct {
	URL   string
	Title string
	Story []Item
}

type Item interface{}
