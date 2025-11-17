import { onMounted, type Ref, ref } from 'vue';


class BaseUI {

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
}
export default BaseUI