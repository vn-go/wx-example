// BaseUI.ts
import { getCurrentInstance, onMounted, ref, type Ref } from 'vue';
import ApiCaller from './api';
import { findViewPath } from './findViewPath';
import { formatDate as libFormatDate } from './formatDate';
import Modal, { type ModalInstance } from './modal';
import UrlNav from './navigator';
import SessionStore from './sessionStore';
const opener = new Modal("html/modal.html");
const urlNav = new UrlNav();
class APICls {

    private _owner: BaseUI;
    private _caller: ApiCaller;
    constructor(owner: BaseUI) {
        this._owner = owner;
        const self = this;
        this._caller = new ApiCaller(() => {
            return self._owner.sessionStore.get("tk")
        });
    }
    async post(apiEnpoint: string, data?: any) {
        const rs = await this._caller.post(this._owner.getViewPath(), apiEnpoint, data);
        if (!rs.ok) {
            if (rs.status == 401) {
                urlNav.move("/auth/login", `ret=${urlNav.getPathname()}`);
                //urlNav.changeUrl("/auth/login", `ret=${urlNav.getPathname()}`);
                return;
            }
        }
        return rs;

    }
}
export default class BaseUI {
    vueComponent: any;
    rootEle: HTMLElement | null = null;
    private _onMount: (() => void) | null = null;
    viewPath: string;
    remoteCaller: APICls
    sessionStore: SessionStore;

    constructor(viewPath?: string) {
        this.sessionStore = new SessionStore("app-store")
        this.remoteCaller = new APICls(this)
        this.viewPath = viewPath;
        // Lấy instance Vue hiện tại
        this.vueComponent = getCurrentInstance();

        // Lifecycle hook: mounted
        onMounted(() => {
            // Lấy root element an toàn
            this.rootEle =
                (this.vueComponent?.ctx?.$ele as HTMLElement) || // pattern project cũ
                (this.vueComponent?.proxy?.$el as HTMLElement) || // fallback Composition API
                null;
            if (this.rootEle && this.rootEle instanceof HTMLElement) {
                if (this.viewPath) {
                    this.rootEle.setAttribute('view-path', this.viewPath);
                }

            }

            // Gọi callback onMounted nếu có
            if (this._onMount) this._onMount();

            // Gọi onInit() sau mounted
            this.onInit();
        });


    }

    /**
     * Khởi tạo khi component ready
     */
    onInit(): void {
        // override trong class con
    }

    /**
     * Đăng ký callback chạy khi component mounted
     */
    onMounted(fn: () => void) {
        this._onMount = fn;
    }

    /**
     * Tạo reactive ref mới
     */
    newEleRef<T = any>(): Ref<T | null> {
        return ref<T>(null) as any;
    }

    /**
     * Helper format date
     */
    formatDate(dateValue: Date | string | number, format = 'dd/MM/yyyy') {
        return libFormatDate(dateValue, format);
    }

    /**
     * Tạo ref key động nếu chưa có
     */
    refKey(key: string) {
        if (!this[key] || this[key] == null) {
            this[key] = ref();
        }
    }

    /**
     * Lấy component bind theo ref key (chờ DOM render)
     */
    getBindComponent(componentKey: string, maxWaitMs = 1000): Promise<Ref<any> | undefined> {
        return new Promise((resolve) => {
            let elapsed = 0;
            const interval = 50;

            const check = () => {
                const refObj = this[componentKey];
                if (refObj && refObj.value != null) {
                    resolve(refObj);
                } else if (elapsed >= maxWaitMs) {
                    resolve(undefined);
                } else {
                    elapsed += interval;
                    setTimeout(check, interval);
                }
            };
            check();
        });
    }
    newModal(componentPath): ModalInstance {
        if (this.rootEle && this.rootEle instanceof HTMLElement) {
            return opener.load(this.rootEle as HTMLElement, componentPath);
        }

    }
    doClose() {
        const uiId = this.rootEle.getAttribute("ui-id");
        const parentUiId = this.rootEle.parentElement.getAttribute("ui-id");
        if (uiId === parentUiId) {
            this.rootEle.parentElement.parentElement.remove();
        }

    }
    getViewPath() {
        return findViewPath(this.rootEle);
    }
}
