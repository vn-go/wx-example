// src/dynamic-import-helper.js

export const modules = import.meta.glob('./views/**/*.vue')

export function loadComponent(path) {
    const key = `./views${path}.vue` // path = 'system/users'
    const loader = modules[key] //

    if (!loader) {
        return modules['./views/error.vue']()
    }
    return loader()
}