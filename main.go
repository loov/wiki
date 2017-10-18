package main

import (
	"fmt"

	"honnef.co/go/js/dom"
)

func main() {
	wiki := &Wiki{}
	wiki.Pages = append(
		wiki.Pages,
		&Welcome,
		&Second,
	)

	site := &Site{}
	site.Doc = dom.GetWindow().Document()
	site.Wiki = wiki
	site.Node = site.Doc.GetElementByID("app")
	site.Bind()
}

type Site struct {
	Doc  dom.Document
	Node dom.Element

	Wiki   *Wiki
	Stages []*Stage
}

type Stage struct {
	Node dom.Element
	Page *Page
}

func (site *Site) Bind() {
	site.Node.SetInnerHTML("")
	site.Node.AppendChild(site.Render())
}

func (site *Site) Render() dom.DocumentFragment {
	fragment := site.Doc.CreateDocumentFragment()
	fragment.AppendChild(site.RenderHeader())
	fragment.AppendChild(site.RenderPages(site.Wiki.Pages))
	return fragment
}

func (site *Site) RenderHeader() dom.Element {
	eheader := site.Doc.CreateElement("div")
	eheader.Class().Add("header")
	eheader.AppendChild(site.RenderSearch())
	return eheader
}

func (site *Site) RenderSearch() dom.Element {
	esearch := site.Doc.CreateElement("form")
	esearch.Class().Add("search")

	einput := site.Doc.CreateElement("input")
	einput.Class().Add("search-input")
	einput.SetAttribute("type", "search")
	esearch.AppendChild(einput)

	ebutton := site.Doc.CreateElement("button")
	ebutton.SetInnerHTML("Search")
	esearch.AppendChild(ebutton)

	return esearch
}
func (site *Site) RenderPages(pages []*Page) dom.Element {
	epages := site.Doc.CreateElement("div")
	epages.Class().Add("pages")
	for i, page := range pages {
		epages.AppendChild(site.RenderPage(page, i == 0))
	}
	return epages
}

func (site *Site) RenderPage(page *Page, selected bool) dom.Element {
	estage := site.Doc.CreateElement("div")
	estage.Class().Add("stage")
	if selected {
		estage.Class().Add("selected")
	}

	estatus := site.Doc.CreateElement("div")
	estatus.Class().Add("status")
	estage.AppendChild(estatus)

	{
		eedit := site.Doc.CreateElement("div")
		eedit.Class().Add("icon")
		eedit.SetInnerHTML("Edit")
		estatus.AppendChild(eedit)

		eurl := site.Doc.CreateElement("div")
		eurl.Class().Add("url")
		eurl.SetInnerHTML(page.URL)
		estatus.AppendChild(eurl)
	}
	epage := site.Doc.CreateElement("div")
	epage.Class().Add("page")
	estage.AppendChild(epage)

	etitle := site.Doc.CreateElement("div")
	etitle.Class().Add("title")
	etitle.SetInnerHTML(page.Title)
	epage.AppendChild(etitle)

	estory := site.Doc.CreateElement("div")
	estory.Class().Add("story")
	for _, item := range page.Story {
		estory.AppendChild(site.RenderItem(item))
	}
	epage.AppendChild(estory)

	return estage
}

func (site *Site) RenderItem(item Item) dom.Element {
	itemel := site.Doc.CreateElement("div")
	itemel.Class().Add("item")

	switch item := item.(type) {
	case *Paragraph:
		itemel.Class().Add("paragraph")
		p := site.Doc.CreateElement("p")
		p.SetInnerHTML(item.Text)
		itemel.AppendChild(p)
	default:
		itemel.Class().Add("missing")
		itemel.SetInnerHTML(fmt.Sprintf("%+v", item))
	}

	return itemel
}

type Wiki struct {
	Pages []*Page
}

type Page struct {
	URL   string
	Title string
	Story []Item
}

type Item interface{}

type Paragraph struct {
	Text string
}

var Welcome = Page{
	URL:   "http://localhost:8000/welcome",
	Title: "Welcome",
	Story: []Item{
		&Paragraph{"Lorem [[ipsum]] dolor sit amet, consectetur adipisicing elit. Optio rerum nam architecto [[mollitia]] unde officia tempora reiciendis omnis quas, expedita, quia ad culpa. Porro explicabo temporibus, sunt officia vel corporis."},
		&Paragraph{"Lorem [[ipsum]] dolor sit amet, consectetur adipisicing elit. Quae dolores consequatur officia [[at]] rem voluptatibus hic eos quo commodi minima? Ipsum commodi quae dolorum eum repudiandae provident saepe reiciendis vitae."},
		&Paragraph{"Lorem [[ipsum]] dolor sit amet, consectetur adipisicing elit. Voluptate porro obcaecati numquam [[asperiores]] dignissimos accusantium veniam, dolore, illo doloremque dolorum, qui aliquam repudiandae nihil neque autem. Quaerat mollitia error molestias."},
		&Paragraph{"Lorem [[ipsum]] dolor sit amet, consectetur adipisicing elit. Minima natus temporibus [[minus]] officia, repudiandae asperiores, nemo cum fugit rem repellat sequi iusto vero, explicabo corporis! Perspiciatis sequi, consectetur eveniet voluptate!"},
		&Paragraph{"Lorem [[ipsum]] dolor sit amet, consectetur adipisicing elit. Optio rerum nam architecto [[mollitia]] unde officia tempora reiciendis omnis quas, expedita, quia ad culpa. Porro explicabo temporibus, sunt officia vel corporis."},
		&Paragraph{"Lorem [[ipsum]] dolor sit amet, consectetur adipisicing elit. Quae dolores consequatur officia [[at]] rem voluptatibus hic eos quo commodi minima? Ipsum commodi quae dolorum eum repudiandae provident saepe reiciendis vitae."},
		&Paragraph{"Lorem [[ipsum]] dolor sit amet, consectetur adipisicing elit. Voluptate porro obcaecati numquam [[asperiores]] dignissimos accusantium veniam, dolore, illo doloremque dolorum, qui aliquam repudiandae nihil neque autem. Quaerat mollitia error molestias."},
		&Paragraph{"Lorem [[ipsum]] dolor sit amet, consectetur adipisicing elit. Minima natus temporibus [[minus]] officia, repudiandae asperiores, nemo cum fugit rem repellat sequi iusto vero, explicabo corporis! Perspiciatis sequi, consectetur eveniet voluptate!"},
	},
}

var Second = Page{
	URL:   "http://localhost:8000/second",
	Title: "Second",
	Story: []Item{
		&Paragraph{"Lorem ipsum dolor sit amet, consectetur adipisicing elit. Minus dolor ab aliquam sunt sit, eaque animi ut recusandae! Vitae ad rem eaque accusantium ex distinctio temporibus quo tempore? Vero, deserunt."},
		&Paragraph{"Lorem ipsum dolor sit amet, consectetur adipisicing elit. In consectetur nesciunt natus eos quos voluptates debitis consequuntur! Maxime quod libero ipsa sed, nihil at, quam rem consectetur corrupti sequi corporis."},
		&Paragraph{"Lorem ipsum dolor sit amet, consectetur adipisicing elit. Neque doloremque eum amet atque vel consectetur veritatis dolore necessitatibus nostrum voluptatem facilis animi nulla dolorem tempore illum molestias cupiditate, pariatur, enim."},
		&Paragraph{"Lorem ipsum dolor sit amet, consectetur adipisicing elit. Omnis placeat facilis enim aliquam aliquid voluptatum sequi accusantium repellat nihil nemo quasi excepturi mollitia explicabo quas, minus, corporis voluptatem pariatur in?"},
	},
}
