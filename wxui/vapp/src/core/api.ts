const API_BASE = (import.meta as any).env.VITE_API_BASE_URL;

export type ApiResult = {
    ok?: boolean;
    error?: {

        status: number;
        statusText: string;
        data: string;
    };
    data?: any;
};

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
    async post(apiEndpoint: string, data?: any): Promise<ApiResult> {
        const url = this.buildUrl(apiEndpoint);
        const headers: Record<string, string> = { "Content-Type": "application/json" };
        if (this._accessToken) headers["Authorization"] = `Bearer ${this._accessToken()}`;

        return this._fetch(url, {
            method: "POST",
            headers,
            body: data ? JSON.stringify(data) : null,
        });
    }

    /** POST form-urlencoded */
    async formPost(apiEndpoint: string, data?: Record<string, any>): Promise<ApiResult> {
        const url = this.buildUrl(apiEndpoint);
        const headers: Record<string, string> = { "Content-Type": "application/x-www-form-urlencoded" };
        if (this._accessToken) headers["Authorization"] = `Bearer ${this._accessToken()}`;

        const formData = new URLSearchParams(data || {}).toString();

        return this._fetch(url, {
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
        this._onDial();
        const ret: ApiResult = {};
        try {
            const res = await fetch(url, options);

            if (!res.ok) {
                ret.error = {

                    status: res.status,
                    statusText: res.statusText,
                    data: await res.text(),
                };
                return ret;
            }

            // Try parse JSON, nếu empty body thì data = null
            try {
                ret.data = res.status !== 204 ? await res.json() : null;
            } catch (err: any) {
                ret.error = {

                    status: res.status,
                    statusText: "Invalid JSON",
                    data: err.message || "Failed to parse JSON",
                };
            }

            return ret;
        } catch (err: any) {
            ret.error = {

                status: 0,
                statusText: "Network Error",
                data: err.message || "Unknown network error",
            };
            return ret;
        } finally {
            this._onFinished();
        }
    }
}

export default ApiCaller;
