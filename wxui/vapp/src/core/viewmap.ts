export default function getViewMap() {
    const viewMaps = (import.meta as any).glob('./../views/**/*.vue', { eager: false });
    return viewMaps
};
import { defineAsyncComponent, markRaw } from 'vue';
export async function loadViews(viewPath?: string, errorView?: string) {
    const viewsData = getViewMap();
    if (viewPath.startsWith('/')) {
        viewPath = viewPath.substring(1, viewPath.length);
    }
    if (!viewsData[`../views/${viewPath}.vue`]) {
        return markRaw(defineAsyncComponent(viewsData[`../views/${errorView}.vue`]));
    }
    return markRaw(defineAsyncComponent(viewsData[`../views/${viewPath}.vue`]));
}