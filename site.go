package main

import (
	"fmt"

	"honnef.co/go/js/dom"
)

type Site struct {
	Node      dom.Element
	StageList dom.Element

	Wiki *Wiki
}

func (site *Site) Init() {
	site.Node.SetInnerHTML("")
	site.Node.AppendChild(site.RenderHeader())

	site.StageList = h("div", "stages")
	site.Node.AppendChild(site.StageList)
}

func (site *Site) RenderHeader() dom.Element {
	eheader := h("div", "header")
	eheader.AppendChild(site.RenderSearch())
	return eheader
}

func (site *Site) RenderSearch() dom.Element {
	esearch := h("form", "search")

	einput := h("input", "search-input")
	einput.SetAttribute("type", "search")
	esearch.AppendChild(einput)

	ebutton := h("button", "")
	ebutton.SetInnerHTML("Search")
	esearch.AppendChild(ebutton)

	return esearch
}
func (site *Site) UpdateStages() {
	site.StageList.SetInnerHTML("")
	for i, page := range site.Wiki.Pages {
		site.StageList.AppendChild(site.RenderPage(page, i == 0))
	}
}

func (site *Site) RenderPage(page *Page, selected bool) dom.Element {
	estage := h("div", "stage")
	if selected {
		estage.Class().Add("selected")
	}

	estatus := h("div", "status")
	estage.AppendChild(estatus)

	{
		eedit := h("div", "icon")
		eedit.SetInnerHTML("Edit")
		estatus.AppendChild(eedit)

		eurl := h("div", "url")
		eurl.SetInnerHTML(page.URL)
		estatus.AppendChild(eurl)
	}
	epage := h("div", "page")
	AttachOverflowIndicator(epage)
	estage.AppendChild(epage)

	etitle := h("div", "title")
	etitle.SetInnerHTML(page.Title)
	epage.AppendChild(etitle)

	estory := h("div", "story")
	for _, item := range page.Story {
		estory.AppendChild(site.RenderItem(item))
	}
	epage.AppendChild(estory)

	return estage
}

func (site *Site) RenderItem(item Item) dom.Element {
	itemel := h("div", "item")

	switch item := item.(type) {
	case *Paragraph:
		itemel.Class().Add("paragraph")
		p := h("p", "")
		(&Parser{
			Begin: func() { p = h("p", "") },
			End:   func() { itemel.AppendChild(p); p = nil },
			Text: func(s string) {
				p.AppendChild(text(s))
			},
			Link: func(spec string) {
				link := h("a", "")
				link.SetAttribute("href", spec)
				link.SetTextContent(spec)
				p.AppendChild(link)
			},
		}).Run(item.Text)
	default:
		itemel.Class().Add("missing")
		itemel.SetInnerHTML(fmt.Sprintf("%+v", item))
	}

	return itemel
}
