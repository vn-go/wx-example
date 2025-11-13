import { onDestroy, onMount } from "svelte";
export default class BaseComponent {
    private eventHandlers: Map<string, (...args: any[]) => void>;

    constructor() {
        this.eventHandlers = new Map();
    }

    on(eventName: string, handler: (sender: any, ...args: any[]) => void) {
        this.eventHandlers.set(eventName, handler);
    }

    off(eventName: string) {
        this.eventHandlers.delete(eventName);
    }

    emit(eventName: string, ...args: any[]) {
        const handler = this.eventHandlers.get(eventName);
        if (handler) handler(this, ...args);
    }
    onMount(callback: () => void) {
        onMount(callback);
    }
    onDestroy(callback: () => void) {
        onDestroy(callback);
    }
}