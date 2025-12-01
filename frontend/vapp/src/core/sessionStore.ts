
export default class SessionStore {
    storeKey: string;

    constructor(key: string) {
        this.storeKey = key;
    }
    get(key: string) {
        let strContent = sessionStorage.getItem(this.storeKey);
        let data = {}
        if (strContent) {
            data = JSON.parse(strContent);
        }
        return data[key.toLowerCase()];

    }
    set(key: string, val: any) {
        let strContent = sessionStorage.getItem(this.storeKey);
        let data = {}
        if (strContent) {
            data = JSON.parse(strContent);
        }
        data[key.toLowerCase()] = val;
        sessionStorage.setItem(this.storeKey, JSON.stringify(data));
    }
}