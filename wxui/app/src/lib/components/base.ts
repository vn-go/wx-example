import { apiCall } from '@lib/utils/apis';
import { onMount } from 'svelte';
export function module(moduleName: string) {
    return function <T extends { new(...args: any[]): {} }>(constructor: T) {
        // Tạo subclass wrapper để chèn logic
        return class extends constructor {
            viewPath = `${moduleName}/${constructor.name}`;
        };
    };
}
export class UIForm {
    viewPath: string;
    OnMounth(callback: () => void) {
        onMount(callback);
    }
    async PostData(apiEnpoint: string, data: any): Promise<any> {
        let response = await apiCall.post(apiEnpoint, data, {
            'View-Path': this.viewPath
        });
        return response.data as any;
    }
    constructor(module: string) {

        this.viewPath = `${module}/${this.constructor.name}`;
    }
}