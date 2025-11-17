import { ref, type Ref } from 'vue';
import UrlNav from './navigator';
import BaseUI from "./ui";
import getViewMap, { loadViews } from './viewmap';
const libs = {
    BaseUI: BaseUI,
    newDOMRef: () => {
        let ret = ref<HTMLElement | null>(null);
        return ret as any;
    },
    urlNav: new UrlNav(),
    getViewMap: getViewMap,
    loadViews: loadViews,
    newRef: <T>(val?: T): Ref<T> => {
        let ret = ref(val);
        return ret as any;
    }
}
export default libs;