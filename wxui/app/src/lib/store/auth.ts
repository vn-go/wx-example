// // src/stores/auth.ts
// import { writable } from 'svelte/store';

// // store chính để lưu access token
// export const accessToken = writable<string | null>(null);

// // nếu muốn tiện hơn, có thể thêm helper:
// export function setAccessToken(token: string | null) {
//     accessToken.set(token);
// }

// export function clearAccessToken() {
//     accessToken.set(null);
// }
// src/stores/auth.ts
import { browser } from '$app/environment'; // nếu dùng SvelteKit
import { get, writable } from 'svelte/store';
// nếu là Svelte standalone (không SvelteKit), bạn dùng typeof window !== 'undefined'

export const accessToken = writable<string | null>(null);

// Khi client mount, sync lại token từ sessionStorage
if (browser && typeof sessionStorage !== 'undefined') {
    const stored = sessionStorage.getItem('token');
    if (stored) {
        accessToken.set(stored);
    }
}

// Lắng nghe thay đổi và lưu lại vào sessionStorage
accessToken.subscribe((value) => {
    if (!browser || typeof sessionStorage === 'undefined') return;
    if (value) sessionStorage.setItem('token', value);
    else sessionStorage.removeItem('token');
});

// Helper
export function setAccessToken(token: string | null) {
    accessToken.set(token);
}
export function getAccessToken() {
    return get(accessToken);
}
export function clearAccessToken() {
    accessToken.set(null);
}
