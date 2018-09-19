package jira

import "time"

type Page struct {
	Version  int       `json:"version"`
	Title    string    `json:"title"`
	Synopsis string    `json:"synopsis,omitempty"`
	Modified time.Time `json:"modified,omitempty"`
	Story    []Item    `json:"story,omitempty"`
}

func (page *Page) Add(items ...Item) {
	page.Story = append(page.Story, items...)
}

type Item map[string]interface{}

// String returns a string value from key
func (item Item) String(key string) string {
	if s, ok := item[key].(string); ok {
		return s
	}
	return ""
}

// Type returns the item `type`
func (item Item) Type() string { return item.String("type") }

// ID returns the `item` identificator
func (item Item) ID() string { return item.String("id") }
