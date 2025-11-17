import { slideFromRightToLeft } from './animation';
import login from './api.login';
import { renderDynamicComponent } from './dynamic-render';
import { slidePageIn } from './slidePageIn';
import { slidePageTransition } from './slidePageTransition';
// import { slideReplaceRight } from './slideReplaceRight';
type PageChangeHandlerArgs = {
    PagePath: string | undefined;
    Query: Record<string, any>;
};

type PageChangeHandler = (data: PageChangeHandlerArgs) => void;

export class Application {

    sessionKey: string;
    getSessionValue<T>(key: string): T | undefined {
        let strData = sessionStorage.getItem(this.sessionKey)
        let data = {}
        if (strData) {
            data = JSON.parse(strData)
        }
        return (data as any)[key]
    }
    setSessionValue<T>(key: string, val: T | undefined) {
        let strData = sessionStorage.getItem(this.sessionKey)
        let data = {}
        if (strData) {
            data = JSON.parse(strData)
        }
        (data as any)[key] = val;
        sessionStorage.setItem(this.sessionKey, JSON.stringify(data))
    }
    private _onUrlChange?: PageChangeHandler;
    private hasRaisedInitialEvent = false;
    pathname: string;

    constructor() {
        this.sessionKey = "application";
        this.hasRaisedInitialEvent = false;
        this.pathname = window.location.pathname;
        if (this.pathname.startsWith('/')) {
            this.pathname = this.pathname.substring(1, this.pathname.length)
        }
        window.addEventListener('popstate', (event) => {
            let nextPath = window.location.pathname;
            if (nextPath.startsWith('/')) {
                nextPath = nextPath.substring(1, window.location.pathname.length);
            }
            if (this.pathname != nextPath) {
                this.pathname = nextPath;

                this.raiseCurrentUrl();
            }

        });


    }
    /**
     * Parses a query string into an object with camelCase keys.
     * Uses kebab-to-camel conversion on keys and decodes values.
     *
     * @param queryString - The input query string (e.g., "abc-cde=1&fx-ycc=hello")
     * @returns An object with camelCase keys (e.g., { abcCde: 1, fxYcc: 'hello' })
     *
     * @example
     *   parseQueryToCamelObject("abc-cde=1&fx-ycc=hello")  // returns { abcCde: 1, fxYcc: 'hello' }
     *   parseQueryToCamelObject("")                       // returns {}
     *   parseQueryToCamelObject("key=value")              // returns { key: 'value' }
     */
    parseQueryToCamelObject(queryString: string): Record<string, any> {
        const result: Record<string, any> = {};

        if (!queryString) return result;
        queryString = queryString.split('?')[1];
        if (!queryString) return result;
        // Split into key-value pairs
        const pairs = queryString.split('&');

        for (const pair of pairs) {
            const [key, value = ''] = pair.split('=', 2);
            if (key) {
                // Convert kebab-case key to camelCase
                const camelKey = this.kebabToCamelCase(key);
                // Decode value (handle %20, + etc.)
                const decodedValue = decodeURIComponent(value.replace(/\+/g, ' '));
                // Parse numbers if possible, else keep as string
                const parsedValue = /^\d+(\.\d+)?$/.test(decodedValue)
                    ? Number(decodedValue)
                    : decodedValue;
                result[camelKey] = parsedValue;
            }
        }

        return result;
    }
    /**
     * Converts a kebab-case string to camelCase.
     * 
     * @param kebabKey - The input string in kebab-case (e.g., "page-index", "user-profile-url")
     * @returns The string in camelCase (e.g., "pageIndex", "userProfileUrl")
     * 
     * @example
     *   kebabToCamelCase("page-index")        // returns "pageIndex"
     *   kebabToCamelCase("user-profile-url")  // returns "userProfileUrl"
     *   kebabToCamelCase("API-KEY")           // returns "apiKey"
     *   kebabToCamelCase("")                  // returns ""
     */
    kebabToCamelCase(kebabKey: string): string {
        if (!kebabKey) return '';

        return kebabKey
            .toLowerCase() // Ensure consistent base
            .split('-')
            .filter(part => part.length > 0) // Remove empty parts
            .map((part, index) =>
                index === 0
                    ? part // First part stays lowercase
                    : part.charAt(0).toUpperCase() + part.slice(1) // Capitalize rest
            )
            .join('');
    }
    /**
     * Converts a camelCase string to kebab-case.
     * 
     * @param camelCase - The input string in camelCase (e.g., "pageIndex", "userProfileUrl")
     * @returns The string in kebab-case (e.g., "page-index", "user-profile-url")
     * 
     * @example
     *   camelToKebabCase("pageIndex")        // returns "page-index"
     *   camelToKebabCase("userProfileUrl")   // returns "user-profile-url"
     *   camelToKebabCase("APIKey")           // returns "api-key"
     *   camelToKebabCase("")                 // returns ""
     */
    camelToKebabCase(camelCase: string): string {
        if (!camelCase) return '';

        return camelCase
            .replace(/[A-Z]/g, (match, offset) =>
                (offset > 0 ? '-' : '') + match.toLowerCase()
            )
            .replace(/^-+/, ''); // Remove leading hyphens if any
    }
    /** Đăng ký listener cho sự kiện thay đổi URL */
    onUrlChange(fn: PageChangeHandler): void {
        this._onUrlChange = fn;

        // Chỉ gọi 1 lần khi khởi tạo (nếu chưa gọi)
        if (!this.hasRaisedInitialEvent) {
            this.hasRaisedInitialEvent = true;
            this.raiseCurrentUrl();
        }
    }

    /** Gọi handler với URL hiện tại */
    raiseCurrentUrl(): void {
        if (!this._onUrlChange) return;

        const data: PageChangeHandlerArgs = {
            PagePath: window.location.pathname || undefined,
            Query: this.parseQueryToCamelObject(window.location.search),
        };

        this._onUrlChange(data);
    }

    /** Phân tích query string thành Map<string, string> */
    private parseSearchParams(search: string): Map<string, string> {
        const map = new Map<string, string>();

        if (!search || search === '?' || search.length <= 1) {
            return map;
        }

        // Bỏ dấu ? đầu tiên
        const queryString = search.startsWith('?') ? search.slice(1) : search;
        const pairs = queryString.split('&');

        for (const pair of pairs) {
            const [key, value = ''] = pair.split('=', 2);
            if (key) {
                // Giải mã URL và chuẩn hóa key
                const decodedKey = decodeURIComponent(key).toLowerCase();
                const decodedValue = decodeURIComponent(value.replace(/\+/g, ' '));
                map.set(this.kebabToCamelCase(decodedKey), decodedValue);
            }
        }

        return map;
    }
    async loginAsync(username: string, password: string) {
        return await login(username, password);
    }
    async renderAsync(
        virtualPath: string,
        target: HTMLElement | undefined,
        props: Record<string, any> = {}) {
        let ret = await renderDynamicComponent(virtualPath, target, props)
        return ret;
    }
    async loadErrorPage(target: HTMLElement | undefined) {
        return await this.renderAsync("error", target);
    }
    async slideFromRightToLeft(fromEle: HTMLElement | undefined,
        toEle: HTMLElement | undefined,
        duration: number = 400) {
        slideFromRightToLeft(fromEle, toEle);
    }
    async slidePageTransition(
        newEle: HTMLElement | undefined,
        container: HTMLElement | undefined,
        duration: number = 500) {
        await slidePageTransition(newEle, container, duration)
    }
    async slidePageIn(newEle: HTMLElement | undefined,
        container: HTMLElement | undefined,
        duration: number = 500) {
        await slidePageIn(newEle, container, duration)
    }
    slideReplaceRight(newEle: any, container: any) {
        slidePageTransition(newEle, container)
    }
}
const app = new Application();
export default app;
