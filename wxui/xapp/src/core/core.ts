// application.ts


import { onDestroy, onMount } from 'svelte';
import { writable, type Writable } from 'svelte/store';
import app, { Application } from './core.Application';
import { findBySlotName } from './dom.find-slot';
// Export instance
export const application = new Application();

// Export class BaseUi (nếu cần dùng sau)
export class BaseUi {
    app: Application;
    Element: HTMLElement | undefined;
    private _onParamChange?: (data: Record<string, any>) => void;
    private _hasRaiseEvent: boolean;
    private _onMount: (() => void) | undefined;
    private _onDismount: (() => void) | undefined;
    pageId: any;
    constructor() {
        this.app = app;
        this._hasRaiseEvent = false;
        window.addEventListener('popstate', (event) => {
            console.log(event);
        });
        onMount(() => {

            if (this._onMount) {
                this._onMount();
            }
            this.onInit();
            if (this.Element) {
                (this.Element as any).__instance = this;

                console.log(this.Element.parentNode);
                this.pageId = (this.Element.parentNode as any).getAttribute("page-path")
                this.Element.setAttribute("view-id", this.pageId);
            }
        })
        onDestroy(() => {
            if (this._onDismount) {
                this._onDismount();
            }
            this.onDestroy();
        })

    }
    onInit() { }
    onDestroy() { }
    async findBySlotName(ele: HTMLElement | undefined, slotName: string, timeout?: number) {
        return await findBySlotName(ele, slotName, timeout)
    }
    /**
 * Navigates to a new page path with optional query parameters.
 * Converts object keys to kebab-case and properly encodes values.
 *
 * @param pagePath - The target URL path (e.g., "/products" or "/user/profile")
 * @param data - Optional object containing query parameters (e.g., { pageIndex: 1, pageSize: 20 })
 */
    jumpToUrl(pagePath: string, data?: Record<string, any>): void {
        let search = '';

        // If data is provided, build the query string
        if (data && typeof data === 'object') {
            const params = Object.entries(data)
                .filter(([_, value]) => value !== undefined && value !== null) // Skip undefined/null
                .map(([key, value]) => {
                    // Convert camelCase/PascalCase to kebab-case (e.g., pageIndex → page-index)
                    const kebabKey = key.replace(/[A-Z]/g, (match, offset) =>
                        (offset > 0 ? '-' : '') + match.toLowerCase()
                    );
                    // Encode value and replace spaces with %20 (standard URL encoding)
                    const encodedValue = encodeURIComponent(String(value));
                    return `${kebabKey}=${encodedValue}`;
                });

            if (params.length > 0) {
                search = params.join('&');
            }
        }

        // Build final URL
        const url = search ? `${pagePath}?${search}` : pagePath;

        // Update browser URL without reloading the page
        window.history.pushState({}, '', url);
        if (pagePath != this.app.pathname) {
            this.app.pathname = pagePath;
            if (this.app.pathname.startsWith('/')) {
                this.app.pathname = this.app.pathname.substring(1, this.app.pathname.length)
            }
            this.app.raiseCurrentUrl();
        }
        // Optional: Trigger your URL change handler
        // this.raiseUrlChangeEvent(); // Uncomment if you have an event system
    }
    onParamChange(fn: (data: Record<string, any>) => void) {
        this._onParamChange = fn;


    }
    onMount(fn: () => void) {
        this._onMount = fn;
    }
    onDismount(fn: () => void) {
        this._onDismount = fn;
    }
    makeWritable<T1 extends string | number | symbol, T2>(): Writable<Record<T1, T2>> {
        return writable<Record<T1, T2>>({} as Record<T1, T2>);
    }
}