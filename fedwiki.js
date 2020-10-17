const fedwiki = (function () {
    "use strict";

    const renderer = new md.Renderer();
    const options = {
        gfm: true,
        sanitize: true,
        taskLists: true,
        renderer: renderer,
        linksInNewTab: true,
        breaks: true
    };

    const Loading = "loading";
    const Errored = "errored";
    const Denied = "denied";
    const Missing = "missing";
    const Loaded = "loaded";

    class Context {
        constructor(host) {
            if (!host.startsWith("http://") && !host.startsWith("https://")) {
                host = "http://" + host;
            }
            if (!host.endsWith("/")) {
                host += "/"
            }
            this.host = host;
        }

        createURL(slug) {
            return this.host + slug + ".json";
        }

        open(title, slug, optionalForkedFrom) {
            let context = this;
            let url = slug;
            if (slug.startsWith("http://") || slug.startsWith("https://")) {
                // TODO: this doesn't look right
            } else {
                url = this.createURL(slug);
            }

            return new View(context, title, slug, url, optionalForkedFrom);
        }
    }

    class View {
        constructor(context, title, slug, url, optionalForkedFrom) {
            this.context = context;
            this.slug = slug;
            this.url = url;
            this.forkedFrom = optionalForkedFrom;

            this.title = title;
            this.status = Loading;
            this.stage = null;

            this.error = null;
            this.page = null;
        }

        attach(stage) {
            this.stage = stage;
            this.update();
            this.fetch();
        }
        detach() {
            this.stage = null;
            // TODO: cancel pending requests
        }

        unfork() {
            let view = this;
            if (!view.forkedFrom) {
                return false;
            }

            let forkContext = new Context(view.forkedFrom);
            view.context = forkContext;
            view.forkedFrom = null;
            view.url = forkContext.createURL(view.slug);
            view.update();
            view.fetch();
            return true;
        }

        fetch() {
            let view = this;
            fetch(view.url, {
                method: "GET"
            }).then(response => {
                if (response.status == 404) {
                    if (view.unfork()) {
                        return;
                    }

                    view.status = Missing;
                    view.error = "Page missing";
                    return;
                }
                if (response.status != 200) {
                    view.status = Errored;
                    view.error = response.statusText;
                    return;
                }

                view.status = Loaded;
                return response.json().then(content => {
                    view.page = content;
                    if (view.page.journal) {
                        view.page.journal.forEach((op) => {
                            if (op.type === "fork") {
                                view.forkedFrom = op.site;
                            }
                        });
                    }
                }).catch(error => {
                    view.status = Errored;
                    view.error = error;
                });
            }).catch(error => {
                view.status = Errored;
                view.error = error;
            }).finally(_ => {
                view.update();
            });
        }

        update() {
            if (this.stage == null) return;

            if (this.page && this.page.class != "") {
                this.stage.setTag(this.page.class, true);
            }

            this.stage.setTag("loading", this.status == Loading);
            this.stage.setSlug(h.text(this.url));

            let page = h.div("page");
            switch (this.status) {
                case Loading:
                    page.appendChild(h.fragment(
                        h.h1("story-header", h.text(this.title))
                    ));
                    break;
                case Errored:
                    page.appendChild(h.fragment(
                        h.h1("", h.text("Error")),
                        h.p(this.error)
                    ));
                    break;
                case Denied:
                    page.appendChild(h.fragment(
                        h.h1("", h.text("Access Denied")),
                        h.p(this.error)
                    ));
                    break;
                case Missing:
                    page.appendChild(h.fragment(
                        h.h1("", h.text("Page missing")),
                        h.p(this.error)
                    ));
                    break;
                case Loaded:
                    let view = this;
                    page.appendChild(h.fragment(
                        h.h1("story-header", h.text(this.page.title)),
                        h.div("story", ...this.page.story.map((item) => {
                            return view.renderItem(item);
                        }))
                    ));
                    break;
                default:
                    page.appendChild(h.fragment(
                        h.h1("", h.text("Unknown status")),
                        h.p(this.status)
                    ));
                    break;
            }

            this.stage.setContent(page);
        }

        renderItem(item) {
            function slugify(name) {
                return name.replace(/\s/g, '-').replace(/[^A-Za-z0-9-]/g, '').toLowerCase();
            }

            let el = h.div("item");
            el.classList.add(item.type);
            switch (item.type) {
                case "paragraph":
                    let view = this;
                    let text = item.text.replace(/\[\[([^\]]+)\]\]/gi, (match, name) => {
                        let slug = slugify(name);
                        let href = view.context.host + slug + ".html";
                        return '<a title="view" href="' + href + '" data-slug="' + slug + '" >' + name + '</a>';
                    });

                    // TODO: this is unsafe.
                    let p = h.p();
                    p.innerHTML = text;
                    this.listenClicks(p);

                    el.appendChild(p);
                    break;
                case "markdown":
                    // TODO: sanitize html
                    let mark = h.div("markdown");
                    mark.innerHTML = md(item.text, options);
                    this.listenClicks(mark);
                    el.appendChild(mark);
                    break;
                case "image":
                    let img = h.tag("img", "thumbnail");
                    img.src = item.url;
                    img.title = item.caption;

                    el.appendChild(h.fragment(
                        img,
                        h.p(item.text)
                    ));
                    break;
                case "code":
                    el.appendChild(h.pre("", item.text));
                    break;
                case "factory":
                    el.textContent = item.text;
                    break;
                case "html":
                    // TODO: sanitize
                    el.innerHTML = item.text;
                    break;
                default:
                    el.classList.add("missing");
                    el.innerText = JSON.stringify(item);
            }
            return el;
        }

        listenClicks(el) {
            let view = this;
            let links = el.getElementsByTagName("a");
            for (let i = 0; i < links.length; i++) {
                let link = links[i];
                link.addEventListener("click", event => {
                    view.linkClicked(event);
                });
                link.addEventListener("auxclick", event => {
                    view.linkClicked(event);
                });
            };
        }

        linkClicked(ev) {
            if (this.stage == null) return;

            let target = ev.target;
            ev.stopPropagation();
            ev.preventDefault();

            let slug = target.getAttribute("data-slug");
            let url = target.getAttribute("href");
            if (slug == null) {
                slug = url;
            }

            let child = this.context.open(target.textContent, slug, this.forkedFrom);
            this.stage.open(child, h.isMiddleClick(ev));
        }
    }

    return {
        Context: Context,
        View: View,
    };
})();