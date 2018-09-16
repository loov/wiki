package html

import (
	"strings"

	"github.com/loov/wiki/client"
)

type Context struct {
	Host string
}

func NewContext(host string) *Context {
	return &Context{host}
}

func (context *Context) CreateURL(href string) string {
	return context.Host + href
}

func (context *Context) Open(title, href string) client.View {
	target := context
	url := href
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		target = context
	} else {
		url = context.CreateURL(href)
	}
	return NewView(target, title, url)
}
