

class UrlNav {
    makeQuery(key: string, val: string): string {
        return `${key.toLowerCase()}=${encodeURIComponent(val)}`;
    }
    getPathname() {
        return window.location.pathname.toLowerCase().substring(1);
    }
    addQuery(key, value) {
        const url = new URL(window.location.href);

        // Encode giá trị để tránh ký tự đặc biệt
        const encodedValue = encodeURIComponent(value);

        // Set: nếu có thì update, không có thì add
        url.searchParams.set(key, encodedValue);

        // Cập nhật URL mà không reload trang
        window.history.replaceState({}, "", url.toString());
    }
    getQuery(key) {
        const params = new URLSearchParams(window.location.search);
        const val = params.get(key);
        return val ? decodeURIComponent(val) : null;
    }
    private _onNav: (path?: string, search?: string) => void;
    pathname: string;
    private _init = false;
    init() {
        const self = this;
        self.raiseOnNav();
        if (!this._init) {
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
            this._init = true;
        }
    }
    constructor() {


    }
    raiseOnNav() {
        if (this._onNav) {
            let pathname = window.location.pathname.substring(1);
            this._onNav(pathname, window.location.search.substring(1, window.length))
        }
    }
    onNav(fn: (pathname?: string, search?: string) => void) {
        this._onNav = fn;
    }
    move(pathname?: string, search?: string) {
        if (!pathname) {
            pathname = '/'
        }
        if (pathname[0] != '/') {
            pathname = '/' + pathname
        }
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
    changeUrl(pathname?: string, search?: string) {
        if (search) {
            if ((pathname.toLowerCase() != window.location.pathname.toLowerCase()) ||
                (search.toLowerCase() != window.location.search.substring(1, window.location.search.length).toLocaleLowerCase())) {
                window.history.pushState({}, "", pathname + "?" + search);

            }
        } else {
            if (pathname.toLowerCase() != window.location.pathname.toLowerCase()) {
                window.history.pushState({}, "", pathname);

            }
        }

    }

}

export default UrlNav;