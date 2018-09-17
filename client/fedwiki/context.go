package fedwiki

import (
	"strings"

	"github.com/loov/wiki/client"
)

type Context struct {
	Host string
}

func NewContext(host string) *Context {
	return &Context{Host: host}
}

func (context *Context) CreateURL(slug string) string {
	return context.Host + slug + ".json"
}

func (context *Context) Open(title, slug string) client.View {
	target := context
	url := slug
	if strings.HasPrefix(slug, "http://") || strings.HasPrefix(slug, "https://") {
		target = context
	} else {
		url = context.CreateURL(slug)
	}
	return NewView(target, title, url)
}
