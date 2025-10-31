import { writable } from 'svelte/store';

const storedTheme =
    typeof window !== 'undefined'
        ? localStorage.getItem('theme') || 'light'
        : 'light';

export const theme = writable<'light' | 'dark'>(storedTheme as any);

if (typeof document !== 'undefined') {
    document.documentElement.classList.toggle('dark', storedTheme === 'dark');
}

theme.subscribe((value) => {
    if (typeof document !== 'undefined') {
        document.documentElement.classList.toggle('dark', value === 'dark');
        localStorage.setItem('theme', value);
    }
});
