export class Container {

    constructor() { this.services = new Map() }

    addSingleton(key, factory) {
        this.services.set(key, { factory, instance: null })
    }

    resolve(key) {
        const entry = this.services.get(key)
        if (!entry.instance) entry.instance = entry.factory()
        return entry.instance
    }
}
export default Container