const API_BASE = (import.meta as any).env.VITE_API_BASE_URL;
let _onRequireAuth = undefined;
export function onRequireAuth(fn: () => {}) {
    _onRequireAuth = fn;
}
export type ApiResult = {
    ok?: boolean;
    status?: number;
    error?: {
        statusText: string;
        data: string;
    };
    data?: any;
};
import emitter from "./eventBus";
export class ApiCaller {
    private _accessToken?: () => {} = undefined;

    // Callbacks, default empty
    private _onDial: () => void = () => { };
    private _onFinished: () => void = () => { };

    constructor(onGetAccessToken?: () => {}) {
        this._accessToken = onGetAccessToken;
    }

    /** Set callback khi bắt đầu request */
    onDial(fn: () => void) { this._onDial = fn; }

    /** Set callback khi request kết thúc */
    onFinished(fn: () => void) { this._onFinished = fn; }

    /** POST JSON */
    async post(viewPath: string, apiEndpoint: string, data?: any): Promise<ApiResult> {
        const url = this.buildUrl(apiEndpoint);
        const headers: Record<string, string> = { "Content-Type": "application/json" };
        if (this._accessToken) headers["Authorization"] = `Bearer ${this._accessToken()}`;
        headers["View-Path"] = viewPath;
        return this._fetch(url, {
            method: "POST",
            headers,
            body: data ? JSON.stringify(data) : null,
        });
    }

    /** POST form-urlencoded */
    async formPost(viewPath: string, apiEndpoint: string, data?: Record<string, any>): Promise<ApiResult> {
        const url = this.buildUrl(apiEndpoint);
        const headers: Record<string, string> = { "Content-Type": "application/x-www-form-urlencoded" };
        headers["View-Path"] = viewPath;
        if (this._accessToken) headers["Authorization"] = `Bearer ${this._accessToken()}`;

        const formData = new URLSearchParams(data || {}).toString();

        return await this._fetch(url, {
            method: "POST",
            headers,
            body: formData,
            credentials: "include", // nếu cần cookie
        });
    }

    /** Build final URL, remove double slashes */
    private buildUrl(apiEndpoint: string) {
        return `${API_BASE.replace(/\/$/, "")}/api/${apiEndpoint.replace(/^\//, "")}`;
    }

    /** Private fetch handler */
    private async _fetch(url: string, options: RequestInit): Promise<ApiResult> {

        emitter.emit("on-api-dial");
        const ret: ApiResult = {};
        try {
            const res = await fetch(url, options);
            if (res.status == 401) {
                if (this._accessToken) {
                    emitter.emit('require-login', {});
                    ret.status = res.status;
                    ret.error = {


                        statusText: res.statusText,
                        data: await res.text(),
                    };
                    return ret;
                } else {
                    ret.status = res.status;
                    ret.error = {


                        statusText: res.statusText,
                        data: await res.text(),
                    };
                    return ret;
                }

            }
            if (!res.ok) {

                ret.status = res.status;
                ret.error = {


                    statusText: res.statusText,
                    data: await res.text(),
                };
                return ret;
            }

            // Try parse JSON, nếu empty body thì data = null
            try {
                ret.status = res.status;
                ret.ok = true;
                ret.data = res.status !== 204 ? await res.json() : null;
            } catch (err: any) {
                ret.status = res.status;
                ret.error = {


                    statusText: "Invalid JSON",
                    data: err.message || "Failed to parse JSON",
                };
            }
            finally {
                emitter.emit("on-api-complete");
            }

            return ret;
        } catch (err: any) {
            ret.status = 0;
            ret.error = {


                statusText: "Network Error",
                data: err.message || "Unknown network error",
            };
            return ret;
        } finally {
            emitter.emit("on-api-complete");
        }
    }
}

export default ApiCaller;
