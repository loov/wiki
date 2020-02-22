"use strict";

class h {
    static tag(tag, classNames, ...children) {
        const el = document.createElement(tag);
        if (classNames != null && classNames !== "") {
            classNames.split(/\s+/).forEach(className => {
                el.classList.add(className);
            });
        }
        children.forEach(child => {
            el.appendChild(child);
        });
        return el;
    }
    static text(text) {
        return document.createTextNode(text);
    }
    static p(...texts) {
        const el = document.createElement("p");
        texts.forEach(text => {
            el.appendChild(h.text(text));
        })
        return el;
    }
    static pre(classNames, text) {
        return h.tag("div", classNames, h.text(text));
    }
    static div(classNames, ...children) {
        return h.tag("div", classNames, ...children);
    }
    static h1(classNames, ...children) {
        return h.tag("h1", classNames, ...children);
    }
    static h2(classNames, ...children) {
        return h.tag("h2", classNames, ...children);
    }
    static h3(classNames, ...children) {
        return h.tag("h3", classNames, ...children);
    }
    static form(classNames, ...children) {
        return h.tag("form", classNames, ...children);
    }
    static button(classNames, ...children) {
        return h.tag("button", classNames, ...children);
    }
    static input(classNames, ...children) {
        return h.tag("input", classNames, ...children);
    }
    static fragment(...children) {
        let el = document.createDocumentFragment();
        children.forEach(child => {
            el.appendChild(child);
        });
        return el;
    }

    static isMiddleClick(ev) {
        if (ev instanceof MouseEvent) {
            return ev.button == 1;
        }
        return false;
    }

    static attachOverflowIndicator(elem) {
        let token = 0;

        function update() {
            token = 0;

            let scrollTop = elem.getAttribute("scrollTop");
            let clientHeight = elem.getAttribute("clientHeight");
            let scrollHeight = elem.getAttribute("scrollHeight");

            elem.classList.toggle("overflow-top", scrollTop > 0);
            elem.classList.toggle("overflow-bottom", scrollTop + clientHeight < scrollHeight);
        };

        function handleFrame() {
            window.cancelAnimationFrame(token);
            token = window.requestAnimationFrame(update);
        };

        elem.addEventListener("scroll", {
            passive: true
        }, handleFrame);
        handleFrame();
    }

    static scrollIntoView(el) {
        el.scrollIntoView({
            behavior: "smooth",
            block: "center",
            inline: "center"
        })
    }
}