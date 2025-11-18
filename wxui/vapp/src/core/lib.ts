import { reactive, ref, type Reactive, type Ref } from 'vue';
import login from './api.login';
import { post } from './apiPost';
import { getAppMenuData } from './appMenuData';
import UrlNav from './navigator';
import SessionStore from './sessionStore';
import BaseUI from "./ui";
import getViewMap, { loadViews } from './viewmap';
const libs = {
    BaseUI: BaseUI,
    _afterLogin: undefined,
    newDOMRef: () => {
        let ret = ref<HTMLElement | null>(null);
        return ret as any;
    },
    urlNav: new UrlNav(),
    getViewMap: getViewMap,
    loadViews: async (viewPath?: string, errorView?: string) => {
        return await loadViews(viewPath, errorView);
    },
    newRef: <T>(val?: T): Ref<T> => {
        let ret = ref(val);
        return ret as any;
    },
    getAppMenuData: getAppMenuData,
    sessionStore: new SessionStore("app-store"),
    login: login,
    newReactive: <T>(val: T): Reactive<T> => {
        return reactive(val as any) as any;
    },
    api: {
        post: async (apiEndpoint: string, data?: any) => {
            try {
                return await post(apiEndpoint, data, libs.sessionStore.get("tk"))
            } catch (error) {
                console.log(error);
            }

        }
    },
    raiseAfterLogin: async () => {
        debugger;
        await libs._afterLogin();
    },

    onAfterLogin: (fn: () => void) => {
        libs._afterLogin = fn;
    }
}
export default libs;