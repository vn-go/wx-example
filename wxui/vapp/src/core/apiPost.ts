const API_BASE = (import.meta as any).env.VITE_API_BASE_URL
type ApiResult = {
    error?: {
        ok: boolean; status: number;
        statusText: string; data: string;

    };

    data?: any;

}
export async function post(apiEndpoint: string, data?: any, accessToken?: string): Promise<ApiResult> {
    // Build the final URL (remove double slashes)
    const url = `${API_BASE.replace(/\/$/, "")}/api/${apiEndpoint.replace(/^\//, "")}`;
    let ret: ApiResult = {

    }
    try {
        // Prepare request headers
        const headers: Record<string, string> = {
            "Content-Type": "application/json"
        };

        // Attach access token if provided
        if (accessToken) {
            headers["Authorization"] = `Bearer ${accessToken}`;
        }

        // Send POST request
        const res = await fetch(url, {
            method: "POST",
            headers,
            body: data ? JSON.stringify(data) : null
        });

        // Non-2xx status handling
        if (!res.ok) {
            ret.error = {
                ok: false,
                status: res.status,                 // 0 = network error
                statusText: res.statusText,
                data: res.statusText
            };
            return ret;
        }
        ret.data = await res.json();
        // Parse JSON response

        return ret;
    } catch (err: any) {
        ret.error = {
            ok: false,
            status: 0,                 // 0 = network error
            statusText: "Network Error",
            data: err.message || "Unknown network error"
        }
        return ret;
    }
}