// src/lib/api/apiFetch.ts
import { accessToken } from '@store/auth';
import { showDialog } from '@store/modalManager'; // bạn cần có dialog store/hàm showDialog
import { get } from 'svelte/store';

export async function apiFetch(input: string, init: RequestInit = {}) {
    const token = get(accessToken);

    const headers = {
        ...(init.headers || {}),
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        'Content-Type': 'application/json'
    };

    const response = await fetch(input, { ...init, headers });

    if (response.status === 401) {
        // Token hết hạn hoặc chưa login → hiển thị form đăng nhập
        showDialog('Login');
        throw new Error('Unauthorized');
    }

    if (!response.ok) {
        const text = await response.text();
        throw new Error(text || `Request failed: ${response.status}`);
    }

    return response.json();
}
