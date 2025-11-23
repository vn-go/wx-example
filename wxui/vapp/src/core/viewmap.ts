export default function getViewMap() {
    const viewMaps = (import.meta as any).glob('./../views/**/*.vue', { eager: false });
    return viewMaps
};
import { defineAsyncComponent, markRaw } from 'vue';
export async function loadViews(viewPath?: string, errorView?: string) {
    debugger;
    const viewsData = getViewMap();
    if (viewPath.startsWith('/')) {
        viewPath = viewPath.substring(1, viewPath.length);
    }
    if (!viewsData[`../views/${viewPath}.vue`]) {
        let cmp = defineAsyncComponent(viewsData[`../views/${errorView}.vue`]);
        console.log(cmp);
        let ret = markRaw(cmp);

        return ret;
    }
    let cmp = defineAsyncComponent(viewsData[`../views/${viewPath}.vue`]);
    console.log(cmp);
    let ret = markRaw(cmp);
    //const v = await cmp["__asyncLoader"]()
    //console.log(v);


    return ret;
}