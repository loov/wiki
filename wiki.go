package main

type Wiki struct {
	Pages []*Page
}

type Page struct {
	URL   string
	Title string
	Story []Item
}

type Item interface{}
