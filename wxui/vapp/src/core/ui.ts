import { onMounted, type Ref, ref } from 'vue';
import { formatDate } from './formatDate';

class BaseUI {
    refKey(key: string) {
        if (!this[key] || this[key] == null) {
            this[key] = ref();
        }
    }
    getBindComponent(componentKey: string): Promise<any> {
        if (this[componentKey] && this[componentKey] != null && this[componentKey]._value) {
            return this[componentKey];
        }

        let self = this;

        // Đặt thời gian chờ tối đa (ví dụ: 10 giây)
        const MAX_WAIT_TIME_MS = 10000;

        return new Promise((resolve, reject) => {
            let self = this;
            let count = 0;
            const find = () => {
                if (count > 5) {// timeout 0.5s 
                    resolve(undefined)
                }
                let isOk = self[componentKey] && self[componentKey].value && self[componentKey].value != null;
                if (isOk) {
                    count++;
                    resolve(self[componentKey])
                } else {
                    setTimeout(find, 100);
                }
            }
            setTimeout(find, 100);

        });
    }

    private _onMount: () => void | null;
    constructor() {
        onMounted(() => {
            if (this._onMount) {
                this._onMount();
            }
            this.onInit();
        })
    }
    onInit() {

    }
    newEleRef(): Ref<any, any> {
        return ref(null);
    }
    onMounted(fn: (() => void | null)) {
        this._onMount = fn;
    }
    formatDate(dateValue: Date | string | number, format = "dd/MM/yyyy") {
        return formatDate(dateValue, format);
    }
}
export default BaseUI