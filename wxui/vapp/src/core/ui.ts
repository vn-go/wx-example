// BaseUI.ts
import { getCurrentInstance, onMounted, ref, type Ref } from 'vue';
import { formatDate as libFormatDate } from './formatDate';

export default class BaseUI {
    vueComponent: any;
    rootEle: HTMLElement | null = null;
    private _onMount: (() => void) | null = null;

    constructor() {
        // Lấy instance Vue hiện tại
        this.vueComponent = getCurrentInstance();

        // Lifecycle hook: mounted
        onMounted(() => {
            // Lấy root element an toàn
            this.rootEle =
                (this.vueComponent?.ctx?.$ele as HTMLElement) || // pattern project cũ
                (this.vueComponent?.proxy?.$el as HTMLElement) || // fallback Composition API
                null;

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
}
