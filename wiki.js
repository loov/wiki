const wiki = (function () {
    "use strict";

    class Search {
        constructor() {
            this.node = h.form("search",
                h.input("search-input"),
                h.button("", h.text("Search"))
            );
        }
    }

    class Lineup {
        constructor() {
            this.node = h.div("lineup");
            this.list = [];
            this.contexts = {}; // map fr
        }

        open(host, title, link) {
            let server = this.contexts[host];
            if (server == null) {
                server = this.contexts[""];
            }

            let view = server.open(title, link);
            let stage = new Stage(this, view);

            this.add(stage);
        }

        add(stage) {
            this.list.push(stage);
            this.node.appendChild(stage.node);
        }

        closeTrailing(target) {
            const i = this.list.indexOf(target);
            if (i < 0) {
                return;
            }

            for (let k = this.list.length - 1; k > i; k--) {
                this.closeIndex(k);
            }
        }

        closeIndex(k) {
            // TODO: don't close stages that are being edited
            let stage = this.list[k];
            this.list.splice(k, 1);
            this.node.removeChild(stage.node);
            stage.close();
        }
    }

    class Stage {
        constructor(lineup, view) {
            this.lineup = lineup;
            this.view = view;

            this.content = h.div("content");
            this.content = h.div("content")
            this.slug = h.div("slug")
            this.buttons = h.div("buttons")
            this.node = h.div("stage",
                h.div("indicator"),
                h.div("status",
                    this.slug,
                    this.buttons,
                ),
                this.content
            );

            h.attachOverflowIndicator(this.content);
            this.view.attach(this);
        }

        setTag(className, state) {
            this.node.classList.toggle(className, state);
        }

        setSlug(node) {
            this.slug.innerHTML = "";
            this.slug.appendChild(node);
        }

        setButtons(node) {
            this.buttons.innerHTML = "";
            this.buttons.appendChild(node);
        }

        setContent(node) {
            this.content.innerHTML = "";
            this.content.appendChild(node);
        }

        open(view, append) {
            let next = new Stage(this.lineup, view);
            if (!append) {
                this.lineup.closeTrailing(this);
            }
            this.lineup.add(next);

            h.scrollIntoView(next.node);
        }

        close() {
            this.view.detach();
        }
    }

    class Client {
        constructor() {
            this.search = new Search();
            this.lineup = new Lineup();

            this.node = h.div("app hide-search",
                h.div("header",
                    this.search.node
                ),
                this.lineup.node
            );
        }
    }

    return {
        Search: Search,
        Lineup: Lineup,
        Stage: Stage,
        Client: Client
    }
})()