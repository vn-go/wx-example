import { onMounted } from "vue";
class UrlNav {
    private _onNav: (path?: string, search?: string) => void;
    pathname: string;
    constructor() {
        const self = this;
        onMounted(() => {
            self.raiseOnNav();
            window.addEventListener('popstate', (event) => {
                let nextPath = window.location.pathname;
                if (nextPath.startsWith('/')) {
                    nextPath = nextPath.substring(1, window.location.pathname.length);
                }
                if (self.pathname != nextPath) {
                    self.pathname = nextPath;

                    self.raiseOnNav();
                }

            });
        })
    }
    raiseOnNav() {
        if (this._onNav) {
            this._onNav(window.location.pathname, window.location.search.substring(1, window.length))
        }
    }
    onNav(fn: (pathname?: string, search?: string) => void) {
        this._onNav = fn;
    }
    move(pathname?: string, search?: string) {
        if (search) {
            if ((pathname.toLowerCase() != window.location.pathname.toLowerCase()) ||
                (search.toLowerCase() != window.location.search.substring(1, window.location.search.length).toLocaleLowerCase())) {
                window.history.pushState({}, "", pathname + "?" + search);
                this.raiseOnNav();
            }
        } else {
            if (pathname.toLowerCase() != window.location.pathname.toLowerCase()) {
                window.history.pushState({}, "", pathname);
                this.raiseOnNav();
            }
        }

    }


}

export default UrlNav;