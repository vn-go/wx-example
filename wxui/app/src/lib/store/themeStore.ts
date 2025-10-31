// src/lib/stores/themeStore.ts
import { writable } from 'svelte/store';

const storedTheme = typeof localStorage !== 'undefined'
    ? localStorage.getItem('theme') || 'light'
    : 'light';

export const theme = writable(storedTheme);

// Khi theme thay đổi, cập nhật class cho <html>
theme.subscribe((value) => {
    if (typeof document !== 'undefined') {
        const root = document.documentElement;
        if (value === 'dark') {
            root.classList.add('dark');
        } else {
            root.classList.remove('dark');
        }
    }
    if (typeof localStorage !== 'undefined') {
        localStorage.setItem('theme', value);
    }
});
