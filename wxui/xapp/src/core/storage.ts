// src/lib/storage.ts

import type { Writable } from 'svelte/store';
import { get, writable } from 'svelte/store';
export function writableSession<T>(key: string, initial: T) {
    const { subscribe, set, update } = writable<T>(initial);

    const saved = sessionStorage.getItem(key);
    if (saved !== null) {
        try {
            set(JSON.parse(saved));
        } catch {
            set(initial);
        }
    }

    return {
        subscribe,
        set: (value: T) => {
            try {
                sessionStorage.setItem(key, JSON.stringify(value));
            } catch (e) {
                console.warn('SessionStorage full or blocked', e);
            }
            set(value);
        },
        update: (fn: any) => update((v) => {
            const next = fn(v);
            try {
                sessionStorage.setItem(key, JSON.stringify(next));
            } catch { }
            return next;
        })
    };
}
/**
 * Đọc giá trị hiện tại từ một sessionStore (writableSession).
 * 
 * @param store - Store được tạo bởi writableSession()
 * @returns Giá trị hiện tại trong store (và đã được đồng bộ từ sessionStorage)
 * 
 * @example
 *   const user = writableSession<User | null>('user', null);
 *   const currentUser = readFromSessionStore(user); // → User | null
 */
export function readFromSessionStore<T>(store: Writable<T>): T {
    
    return get(store);
}