const html = (function () {
    "use strict";

    const Loading = "loading";
    const Errored = "errored";
    const Denied = "denied";
    const Missing = "missing";
    const Loaded = "loaded";

    class Context {
        constructor(host) {
            this.host = host;
        }

        createURL(href) {
            return this.host + href;
        }

        open(title, href) {
            let context = this;
            let url = href;
            if (href.startsWith("http://") || href.startsWith("https://")) {
                // TODO: this doesn't look right
                url = href;
            } else {
                url = this.createURL(href);
            }

            return new View(context, title, url);
        }
    }

    class View {
        constructor(context, title, url) {
            this.context = context;
            this.url = url;

            this.title = title;
            this.status = Loading;
            this.stage = null;

            this.error = null;
            this.content = null;
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

        fetch() {
            let view = this;
            let path = "/proxy?url=" + encodeURIComponent(this.url);
            fetch(path, {
                method: "GET"
            }).then(response => {
                if (response.status == 404) {
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
                view.content = response.responseText;
            }).catch(error => {
                view.status = Errored;
                view.error = error;
            }).finally(_ => {
                view.update();
            });
        }

        update() {
            if (this.stage == null) return;

            this.stage.setTag("loading", this.status == Loading);
            this.stage.setSlug(h.text(this.url));

            let page = h.div("page");
            switch (this.status) {
                case Loading:
                // fallthrough
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
                    page.appendChild(this.render());
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

        render() {
            // TODO: sanitize html
            let mark = h.div("markdown");
            mark.innerHTML = this.content;

            let links = mark.getElementsByTagName("a");
            for (let i = 0; i < links.length; i++) {
                let link = links[i];
                link.addEventListener("click", this.linkClicked);
                link.addEventListener("auxclick", this.linkClicked);
            };

            return mark;
        }

        linkClicked(ev) {
            if (this.view == null) return;

            let target = ev.target;
            ev.stopPropgation();
            ev.preventDefault();

            let slug = target.getAttribute("data-slug");
            let url = target.getAttribute("href");
            if (slug == null) {
                slug = url;
            }

            let child = this.context.open(target.textContent, slug);
            this.stage.open(child, h.isMiddleClick(ev));
        }
    }

    return {
        Context: Context,
        View: View,
    };
})();