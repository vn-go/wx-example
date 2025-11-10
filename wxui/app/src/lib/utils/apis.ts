// src/lib/services/ApiCall.ts
export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE'
import { getAccessToken } from "@store/auth";
import { showDialog } from "@store/dialogStore";
export interface ApiResponse<T = any> {
    ok: boolean
    status: number
    data?: T
    error?: string
}
const API_BASE_URL = import.meta.env.VITE_PUBLIC_API_BASE_URL || 'http://localhost:5000/api';
export class ApiCall {
    constructor(private baseUrl = '') { }

    private async request<T>(
        url: string,
        method: HttpMethod,
        body?: any,
        headers: Record<string, string> = {},

    ): Promise<ApiResponse<T>> {
        try {
            const opts: RequestInit = {
                method,
                headers: {
                    'Accept': 'application/json',
                    'authorization': `Bearer ${getAccessToken()}`,
                    ...headers,
                },
                credentials: "include", // 🔥 rất quan trọng
            }

            if (body && method !== 'GET') {
                // Nếu là FormData thì không set content-type
                if (body instanceof FormData) {
                    opts.body = body
                } else {
                    opts.body = JSON.stringify(body)
                    if (opts.headers) {
                        (opts.headers as any)['Content-Type'] = 'application/json'
                    }

                }
            }

            const res = await fetch(this.baseUrl + url, opts)
            const contentType = res.headers.get('content-type') || ''
            let data: any

            if (contentType.includes('application/json')) {
                data = await res.json()
            } else {
                data = await res.text()
            }
            if (res.status === 401) {
                await showDialog('Login')
                return { ok: false, status: 401, error: 'Unauthorized' }
                // TODO: Handle unauthorized error
            }
            return { ok: res.ok, status: res.status, data }
        } catch (err: any) {

            return { ok: false, status: 0, error: err.message || String(err) }
        }
    }

    async get<T>(url: string, headers: Record<string, string> = {}) {
        return this.request<T>(url, 'GET', undefined, headers)
    }

    async post<T>(apiEnpoint: string, body?: any, headers: Record<string, string> = {}) {
        let url = `${API_BASE_URL}/${apiEnpoint}`
        return this.request<T>(url, 'POST', body, headers)
    }

    async formPost<T>(url: string, form: Record<string, any>) {
        const formData = new FormData()
        Object.entries(form).forEach(([k, v]) => {
            if (v !== undefined && v !== null) formData.append(k, v as any)
        })
        return this.request<T>(url, 'POST', formData)
    }
}
export const apiCall = new ApiCall()