import { getCurrentInstance, reactive, ref, type Reactive } from 'vue';

import ApiCaller from './api';
import login from './api.login';
import { getAppMenuData } from './appMenuData';
import Modal from './modal';
import UrlNav from './navigator';
import SessionStore from './sessionStore';
import BaseUI from "./ui";
import getViewMap, { loadViews } from './viewmap';
const sessionStore = new SessionStore("app-store");



const libs = {
    BaseUI: BaseUI,
    _afterLogin: undefined,
    getCurrentInstance: getCurrentInstance,
    getEle: () => {
        return (getCurrentInstance() as any).$ele;
    },
    newDOMRef: () => {
        let ret = ref<HTMLElement | null>(null);
        return ret as any;
    },
    urlNav: new UrlNav(),
    getViewMap: getViewMap,
    loadViews: async (viewPath?: string, errorView?: string) => {
        return await loadViews(viewPath, errorView);
    },
    newRef: <T>(val?: T): any => {
        let ret = ref(val);
        return ret as any;
    },
    getAppMenuData: getAppMenuData,
    sessionStore: sessionStore,
    login: login,
    newReactive: <T>(val: T): Reactive<T> => {
        return reactive(val as any) as any;
    },
    api: new ApiCaller(() => { return sessionStore.get("tk") }),
    apiPublic: new ApiCaller(),
    raiseAfterLogin: async () => {
        debugger;
        await libs._afterLogin();
    },

    onAfterLogin: (fn: () => void) => {
        libs._afterLogin = fn;
    },
    newModal(rawHtml?: string) {
        return new Modal(rawHtml);
    },
    async showModal(rootEle: HTMLHtmlElement, componentPath: string, data?: any) {
        const opener = new Modal("html/modal.html");
        await opener.load(rootEle, componentPath, data);
    },



}
export default libs;