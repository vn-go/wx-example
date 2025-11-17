export default function getViewMap() {
    const viewMaps = (import.meta as any).glob('./../views/**/*.vue', { eager: false });
    return viewMaps
};
import { defineAsyncComponent } from 'vue';
export function loadViews(viewPath?: string) {
    if (viewPath.startsWith('/')) {
        viewPath = viewPath.substring(1, viewPath.length);
    }
    return defineAsyncComponent(getViewMap()[`../views/${viewPath}.vue`]);
}