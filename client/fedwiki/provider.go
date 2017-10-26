package fedwiki

import "github.com/loov/wiki/client"

type Provider struct{}

func (provider *Provider) Open(title, url string) client.Context {
	return NewContext(title, url)
}
